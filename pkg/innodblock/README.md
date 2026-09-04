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
