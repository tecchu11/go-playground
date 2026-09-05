# InnoDB ロック挙動の実測

`DELETE ... WHERE (job_id, ord) IN ((1,1),(1,2))` が複合PKに対して
ギャップロックを取るかどうかを、実機で測った結果と再現手順。

## 実行方法

```sh
go test ./pkg/innodblock/ -v
```

testcontainers で `mysql:8.0.36` を起動し、`testdata/mysql/conf/my.cnf` を
読ませたうえで全ケースを実行する。Docker が必要。
`performance_schema.data_locks` の参照に `PROCESS` 権限が要るため、
テストは root で接続している。

## 測定環境

| 項目 | 値 |
|---|---|
| version | 8.0.36 |
| innodb_version | 8.0.36 |
| transaction_isolation | REPEATABLE-READ |

READ COMMITTED ではそもそもギャップロックを取らないため、
REPEATABLE READ であることを `requireRepeatableRead` で毎回検証している。

## スキーマ

```sql
CREATE TABLE job_attachment (
  job_id     BIGINT NOT NULL,
  ord        INT    NOT NULL,
  object_key VARCHAR(255) NOT NULL,
  PRIMARY KEY (job_id, ord)
) ENGINE=InnoDB;
```

投入データは `(1,1) (1,2) (1,3) (2,1) (2,2) (3,1)`。
`job_id` を跨いで行が隣接しているため、`job_id = 1` の削除が取るギャップが
`(2,1)` の手前として観測できる。

## 結果：ロックモード

`TestDeleteLockModes`。セッションAでトランザクションを開いたまま
`performance_schema.data_locks` の `LOCK_TYPE = 'RECORD'` を読んだもの。

| # | DELETE | LOCK_MODE / LOCK_DATA | ギャップ |
|---|---|---|---|
| 1 | `WHERE job_id = 1` | `X (1,1)` `X (1,2)` `X (1,3)` `X,GAP (2,1)` | あり |
| 2 | `WHERE id IN (2,3)`（単一列PK） | `X,REC_NOT_GAP 2` `X,REC_NOT_GAP 3` | なし |
| 3 | `WHERE (job_id, ord) IN ((1,1),(1,2))` | `X,REC_NOT_GAP (1,1)` `X,REC_NOT_GAP (1,2)` | なし |
| 4 | `WHERE job_id = 1 AND ord IN (1,2)` | `X,REC_NOT_GAP (1,1)` `X,REC_NOT_GAP (1,2)` | なし |
| 5 | `WHERE (job_id, ord) IN ((1,3),(2,1))` | `X,REC_NOT_GAP (1,3)` `X,REC_NOT_GAP (2,1)` | なし |
| 6 | `WHERE (job_id, ord) IN ((1,1),(1,9))`（`(1,9)` は不在） | `X,REC_NOT_GAP (1,1)` `X,GAP (2,1)` | **あり** |
| 7 | `WHERE job_id = 10`（該当行なし） | `X supremum pseudo-record` | **あり** |

Case 3 と Case 4 は同一。どちらの書き方でもギャップは取らない。

## 結果：実行計画

`TestExplainDeletePlans`。3 つの DELETE すべてが `type=range key=PRIMARY`。
Case 1 と Case 3/4 は実行計画が同じ `range` なのにロックが異なる。
つまり `EXPLAIN` の `type` からロック範囲は判断できず、実測が必要だった。

## 結果：デッドロック再現

`TestDeadlockOnRangeDelete` ほか。添付が0件の job 10 / job 11 に対し、
2セッションが delete → insert する。

| テスト | セッションAの DELETE | 結果 |
|---|---|---|
| `TestDeadlockOnRangeDelete` | `WHERE job_id = 10` | **1213 発生** |
| `TestDeadlockOnPrimaryKeyDeleteOfMissingRow` | `WHERE (job_id, ord) IN ((10,1))` | **1213 発生** |
| `TestNoDeadlockWhenNothingIsDeleted` | DELETE を発行しない | 正常 |
| `TestNoDeadlockOnPrimaryKeyDeleteOfExistingRows` | 実在行のPK指定 | 正常 |

`SHOW ENGINE INNODB STATUS` の `LATEST DETECTED DEADLOCK` でも機序が裏取りできる。

```
*** (1) HOLDS THE LOCK(S):
... index PRIMARY of table `lockdemo`.`job_attachment` trx id 1877 lock_mode X
 0: len 8; hex 73757072656d756d; asc supremum;;
*** (1) WAITING FOR THIS LOCK TO BE GRANTED:
... trx id 1877 lock_mode X insert intention waiting
 0: len 8; hex 73757072656d756d; asc supremum;;
```

両トランザクションが supremum に対する同じ next-key lock を保持し
（ギャップロック同士は競合しないので両方 GRANTED される）、
互いに insert intention を待って閉路になる。

## 判定

**複合PK `(job_id, ord)` を維持してよい。** サロゲートPKへの切り替えは不要。

- Case 3/4 が `X,REC_NOT_GAP` のみ。ギャップは取らない
- Case 5 のとおり `job_id` を跨いでも変わらない。効いているのは
  「タプルが完全なPKであること」であって隣接関係ではない

ただし条件が2つある。

1. **実在する行のPKだけを指定すること。** Case 6 のとおり、存在しない
   タプルを名指しすると、そこにあったはずのギャップにロックが降りる。
   PK等値に書き換えただけでは
   `TestDeadlockOnPrimaryKeyDeleteOfMissingRow` のとおりデッドロックは残る。
   現行の添付を先に読み、実在するPKに対してのみ DELETE を発行する必要がある。
2. **削除対象が0件なら DELETE を発行しないこと。** Case 7 のとおり
   該当行のない範囲削除は supremum に next-key lock を残す。
   これが最も刺さりやすい経路で、`TestDeadlockOnRangeDelete` の原因そのもの。

いずれの構成でも 1213 のリトライ機構は入れておく。
ロック取得順序を揃えても、ギャップを一切取らなくなるわけではない
（PK不在・0件削除の経路が残る）ため、リトライは保険として必要。

## ファイル

| ファイル | 内容 |
|---|---|
| `main_test.go` | testcontainers による MySQL 起動 |
| `lockmode_test.go` | スキーマ定義、`data_locks` 観測、ロックモードと実行計画の検証 |
| `deadlock_test.go` | 2セッションを交互に進めるデッドロック再現 |
| `insertlock_test.go` | INSERT 側のロック（A-1） |
| `sequence_test.go` | 本番シーケンス全体の並行実行（A-2） |
| `guardrail_test.go` | 禁止パターンの逆検証と IN リスト件数（B-1 / B-2） |
| `foreignkey_test.go` | FK 制約下での親子ロック（C-1） |

---

# 追加検証（INSERT 側・本番シーケンス）

## A-1. INSERT 単体のロック

`TestInsertLockModes` / `TestInsertedRowIsStillProtected`。

| # | 文（同一TX内） | 保持ロック |
|---|---|---|
| A-1-a | `DELETE ord IN (1,2,3)` → `INSERT (1,1)(1,2)(1,3)` | `X,REC_NOT_GAP` ×3 |
| A-1-b | 同上 → `INSERT (1,1)..(1,5)` | `X,REC_NOT_GAP` ×3（(1,4)(1,5) は 0 件） |
| A-1-c | `DELETE job_id=3 ord IN (1)` → `INSERT (3,1)(3,2)` | `X,REC_NOT_GAP (3,1)` のみ |
| 対照 | `INSERT (3,2)` のみ | ロックなし |

自TXが挿入した行は `data_locks` に明示ロックを持たない（レコードに TRX_ID を書くだけで、
競合が起きたときに初めて明示ロックへ昇格する）。よって `X,GAP` は出ない。

「ロック行が無い＝無防備」ではないことを並行セッションで確認済み。

| 検証 | 結果 |
|---|---|
| 別セッションが同じギャップの別キー `(1,6)` を INSERT | 成功 |
| 別セッションが同じキー `(1,4)` を INSERT | 1205（重複キー防止） |
| 別セッションがテーブル末尾 `(3,3)` を INSERT | 成功 |

## A-2. 本番シーケンス全体の並行実行

`TestSequenceOnAdjacentJobs` ほか。全件正常終了、1213 なし。

| # | 構成 | 結果 |
|---|---|---|
| A-2-a | job 1 / job 2 を交互に | 各TXが自分の行に `X,REC_NOT_GAP` のみ |
| A-2-b | 添付0件の job 10 / job 11（DELETE を発行しない） | `job_attachment` にロック 0 件 |
| A-2-c | job 3 更新 × 新規 job 作成 | 正常終了 |
| A-2-d | 同一 job 1 を2セッション | 親行で待ち1本、A の COMMIT 後に B が進行 |

## B-1. ガードレールの逆検証

`TestForUpdateOnChildSelectTakesGapLocks`。いずれも期待どおりギャップを取る。

| # | 文 | 保持ロック |
|---|---|---|
| B-1-a | `SELECT ord ... WHERE job_id=1 FOR UPDATE` | `X (1,1)(1,2)(1,3)` + `X,GAP (2,1)` |
| B-1-b | `SELECT MAX(ord) ... FOR UPDATE` | `X (1,3)` + `X,GAP (2,1)` |
| B-1-c | `SELECT id FROM job WHERE id=9999 FOR UPDATE` | `X supremum pseudo-record` |

B-1-c は別セッションの `INSERT INTO job` を実際にブロックする（1205 で確認）。

## B-2. IN リスト件数

`TestLargeInListLockModes`。決め手は件数ではなく**選択率**。

| テーブル形状 | サイズ | type | key_len | ロック |
|---|---|---|---|---|
| job 1 が 304 行中 300 行 | 10 | range | 12 | `X,REC_NOT_GAP` ×10 |
| 同上 | 100〜250 | **ALL** | - | **`X` ×304（全行 + supremum）** |
| job 1 が約 5300 行中 300 行 | 10〜250 | range | 12 | `X,REC_NOT_GAP` ×N |

`eq_range_index_dive_limit` は無関係（2 に下げても `X,REC_NOT_GAP` のまま）。
選択率が悪化して `type=ALL` に落ちるとテーブル全体に next-key lock が乗る。

サイズ 230 / 250 で `X supremum pseudo-record` が 1 件追加されることがあるが、
これはページ単位の supremum。別 job への INSERT を一切ブロックしないことを確認済み
（`TestSupremumLockFromLongInListStaysWithinTheJob`）。

## C-1. FK 制約を張った場合

`TestForeignKeyChildInsertLocksParent` ほか。

- 子への INSERT は親行に `S,REC_NOT_GAP` を取る
- C-1-a（親を先に `FOR UPDATE`）：各TXは自分の親行に `X,REC_NOT_GAP` のみ。S は追加されない
- C-1-b：
  - 別 job 同士 → 正常終了
  - 同一 job・シーケンス順（親 UPDATE が先） → 単純な待ち（1205）
  - 同一 job・子 INSERT が先 → **1213**。S を持ったまま X へ昇格しようとして閉路

FK を張るなら、親を先にロックするか、少なくとも親 UPDATE を子書き込みより前に置くことが必須。

## C-2. 他ユースケースとの競合

`job_attachment` を Repository 以外の経路で読み書きする箇所、および `FOR SHARE` を
使う参照系はいずれも存在しないことを確認済み。追加の並行テストは不要。

これにより、前回挙げた TOCTOU（読み取りと削除の間に他TXが行を消す）の経路は閉じている。
本番シーケンスが先頭で `SELECT job FOR UPDATE` を取り、同一 job の2セッションが
そこで直列化されるため（`TestSequenceOnSameJobSerializes` で実測）。

ただしこれは「Repository 以外の経路が無い」ことに依存する前提であり、
親を `FOR UPDATE` しない経路が1つ増えるだけで崩れる。

1213 のリトライは引き続き必要。`TestLargeInListLockModes` の `type=ALL` 転落
（選択率悪化でテーブル全行に next-key lock）が残っているため。
