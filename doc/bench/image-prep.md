# イメージ準備コストの検証（docker save/load vs podman additionalimagestores）

Testcontainers 実行前のイメージ準備コストについて、以下2方式を GitHub Actions 上で比較するための計測基盤。

- **方式A**: `docker save` した tar を actions/cache に保存し、テスト前に `docker load`
- **方式B**: 展開済みの Podman イメージストアを actions/cache に保存し、`additionalimagestores` で参照

検証したい仮説: 方式Bはレイヤ展開が逐次実行される制約を回避できるため、イメージ準備フェーズの所要時間が短縮される。

## 本リポジトリ適用時の前提差分

計画書の前提と本リポジトリの実態には以下の差がある。

- `.github/workflows/ci.yaml` には `docker save` / `docker load` / `actions/cache` が無く、
  「方式A（現行）」は存在しない。毎回 testcontainers が Docker Hub から pull している。
  そのため方式Aは本検証用に新規実装したものであり、ベースラインは実 CI の再現ではない。
- 対象イメージは `cmd/api/internal/datasource/main_test.go` の `testImage` に合わせて `mysql:8.0.36`。
- `ContainerRequest` の `Networks` 明示は Go ファイル全体で 0 件。
  Docker(`bridge`) / Podman(`podman`) のデフォルトネットワーク名差異は本リポジトリでは問題にならない。
  `bench-podman-smoke.yaml` がこの前提を毎回チェックする。

## 構成

| ファイル | 役割 |
|---|---|
| `.github/image-versions.txt` | 対象イメージのピン留め。キャッシュキーの `hashFiles` 対象 |
| `.github/actions/start-podman-socket/` | rootless podman の Docker 互換ソケット起動と `DOCKER_HOST` 設定 |
| `.github/workflows/bench-prepare-cache.yaml` | 方式A/B 双方のキャッシュ生成（手動 / 週次） |
| `.github/workflows/bench-podman-smoke.yaml` | Phase 1。Podman 上でテストが通るかの疎通確認のみ |
| `.github/workflows/bench-image-prep.yaml` | Phase 0 / Phase 3。A/B 交互実行と集計 |

## 実行手順

1. `bench-prepare-cache` を手動実行してキャッシュを作る。
2. `bench-podman-smoke` を手動実行し、rootless podman 上でテストが完走することを確認する。
   ここで落ちる場合は方式Bを中止する。
3. `bench-image-prep` を `iterations: 5`, `methods: both` で実行する。
   ジョブサマリに中央値 / 最小 / 最大の比較表が出力される。

## 計測区間

| 指標 | 方式A | 方式B |
|---|---|---|
| T0 | 0（docker daemon は常駐済み） | podman ソケット起動 |
| T1 | tar の DL + zstd 伸長 | ストア tar の DL + zstd 伸長 |
| T2 | `docker load` | `sudo tar -x`（T2a）+ 追加ストア参照の確認（T2b） |
| T3 | 最初のコンテナ起動完了まで | 同左 |
| T4 | ジョブ全体 | 同左 |
| S | キャッシュサイズ（ディスク実サイズ / actions/cache 上のサイズ） | 同左 |

T3 を独立して測るのは、方式Bで展開コストがコンテナ起動時に後ろ倒しされる可能性があるため。
T1+T2 だけを見て速くなったと誤認しないための指標。

### 比較の公平性のために固定した条件

- **Go のコンパイルを計測ジョブから完全に排除**。`warm` ジョブで `go test -c` したテストバイナリを
  artifact 経由で各計測ジョブへ配布する。
  `actions/setup-go` の既定キャッシュはキーが `go.sum` ハッシュのみで `ci.yaml` と共有されるうえ、
  `--race` の有無でビルドキャッシュの内容が変わるため、計測ジョブでは `cache: false` にしている。
- **Ryuk は両方式で無効化**（`TESTCONTAINERS_RYUK_DISABLED=true`）。
  Podman 側だけ無効にすると Ryuk コンテナの起動時間分だけ方式Aが不利になる。
- **A, B, A, B ... の交互直列実行**（`max-parallel: 1`）。
  actions/cache のスループット変動を時間帯で相殺する。
- **キャッシュキーは固定**（`.github/image-versions.txt` のハッシュのみ）。毎回ヒットする。

### T3 の測り方

`go test -c` で作ったテストバイナリを `-test.run '^$'` で起動する。
テスト本体は 0 件だが `TestMain` は実行されるため、
コンテナ起動 + マイグレーション + フィクスチャ投入までが測られる。

`TestMain` の context timeout は 1 分（`main_test.go`）。
方式Bで展開がコンテナ起動時に後ろ倒しされた場合、T3 の悪化ではなく panic として現れる。

### イメージ投入失敗の検出

投入に失敗しても testcontainers は Docker Hub から pull してテストを通してしまう。
「速くなったように見えて実は rate limit を消費していた」を避けるため、
T2 直後に以下でジョブを失敗させる。

- 方式A: `docker image inspect mysql:8.0.36`
- 方式B: `podman images` の出力に `docker.io/library/mysql:8.0.36 true`（readonly）があること

## 実行して判明した事実（2026-09-03）

### 1. actions/cache は「展開済みのまま」置けない

`actions/cache` は保存時に必ず `tar --posix -cf cache.tzst --use-compress-program zstdmt` を実行し、
復元時に必ず untar する。ディレクトリをそのままの形で置く手段は無い。

したがって方式Bの仮説「load 相当の処理が発生しない」は、actions/cache を使う限り成立しない。
方式Bでも復元時に展開は発生する。差が出るとすれば

- 方式A: `docker load`（daemon の import API 経由でレイヤを逐次展開し overlay2 を構築）
- 方式B: `tar -x`（zstdmt によるマルチスレッド伸長 + tar 展開）

の実装差であって、展開の有無ではない。計測する価値は残るが、
判定基準の「T2 は消えたが T1 が増えて相殺」に該当する可能性が高い前提で読むこと。

### 2. rootless podman のストアは runner 権限で tar できない

rootless podman のイメージストアには sub-UID にマップされたファイルが含まれる
（例: イメージ内の `mysql:mysql` = ホスト上の uid 100998）。
`actions/cache` は `runner` 権限で tar するため、そのままでは読めない。

```
tar: .../diff/var/log/mysqld.log: Cannot open: Permission denied
tar: Exiting with failure status due to previous errors
##[warning]Failed to save: "/usr/bin/tar" failed with error: exit code 2
##[warning]Cache save failed.
```

**この失敗は warning 止まりでジョブは success になる。** 初回実行ではこれに気づかず、
キャッシュが空のまま計測に進むところだった（run 33781250242）。

### 3. 所有者の平坦化はイメージを壊す（この方向は破棄した）

上記を回避するため `sudo chown -R runner:runner` + `sudo chmod -R u+rX` で
アーカイブ可能にした。保存は成功した（run 33783173149）が、
コンテナが起動しなくなった（run 33784022153、方式B 5回すべて失敗）。

```
[ERROR] [MY-010338] Can't find error-message file '/usr/share/mysql-8.0/errmsg.sys'
[ERROR] [MY-012576] [InnoDB] Unable to create temporary file inside "/tmp"; errno: 13
[ERROR] [MY-013236] The designated data directory /var/lib/mysql/ is unusable.
panic: run mysql: ... wait until ready: container exited with code 1
```

mysql の entrypoint は uid 999 (`mysql`) に降りて mysqld を起動する。
所有者を平坦化するとイメージ内の非 root ユーザー所有ファイルが container-root 所有になり、
uid 999 から読めなくなる（errno 13 = EACCES）。

到達した結論は次のとおり。

- 「展開済みストアをそのまま actions/cache に置く」は**できない**（必ず tar される）
- 「アーカイブできるよう所有者を直す」も**できない**（イメージが壊れる）

そこで方式Bを、**所有者を保ったまま root 権限で tar を作り、復元後に root 権限で展開する**
形に組み替えた。所有者は保たれ、比較としても意味が残る。

| | 方式A | 方式B（組み替え後） |
|---|---|---|
| キャッシュの中身 | `docker save` の tar | `sudo tar` で固めたストアの tar |
| T2（イメージ投入） | `docker load`（daemon の逐次展開） | `sudo tar -x` + 追加ストア参照の確認 |

つまり測っているのは **「daemon 経由の逐次展開」対「tar 展開」** であって、
当初の仮説にあった「展開の有無」ではない。

### 4. 途中で踏んだ罠（再発防止のため記録）

- `podman unshare chown -R 0:0` では mysql 固有層は直ったが、ベースイメージ層の
  `/etc/shadow` `/etc/gshadow` が読めないまま残った（run 33781833195）。
- `tar -cf /dev/null` は検証にならない。GNU tar は出力先が `/dev/null` のとき
  ファイル内容を読まない最適化をするため、読めないファイルがあっても成功する。
- `find ! -readable` はシンボリックリンクのリンク先を辿る。mysql イメージには
  壊れたリンク（`/usr/lib/.build-id/*` → 別レイヤの `lib64/*.so`）が含まれ、
  `-type f -o -type d` で限定しないと誤検出する（run 33782309702）。
- 所有者を `runner` にしても mode `----------`（0000）のファイル
  （Debian の `/etc/shadow-` `/etc/gshadow-`）は読めない。所有者にも DAC は適用される
  （run 33782754135）。

### 5. 計測結果

キャッシュ生成時のサイズ（`bench-prepare-cache`）。

| | 方式A | 方式B |
|---|---|---|
| 生サイズ | 604,794,915 B（`docker save` の tar） | 604,794,915 B（ストア） |
| actions/cache 上のサイズ | 151,791,387 B | 156,379,425 B（+3.0%） |

Phase 1（`bench-podman-smoke` run 33783584258）は success。
rootless podman 上で `go test ./... --race -cover` が 92 秒で完走した。
`Networks` の明示は無く、Ryuk 無効化で問題は出ていない。

A/B 計測（`bench-image-prep` run 33785578983、各方式 5 回、A/B 交互直列）。
単位はミリ秒、`中央値 / 最小 / 最大`。

| 指標 | 方式A (n=5) | 方式B (n=5) |
|---|---|---|
| T0 ソケット起動 | 0 / 0 / 0 | 12069 / 9712 / 17071 |
| T1 キャッシュ復元 | 3039 / 1894 / 3446 | 2955 / 1751 / 4503 |
| T2 イメージ投入 | 12043 / 11069 / 17113 | 1395 / 1045 / 1521 |
| **T1+T2** | **15489** / 13701 / 20152 | **4000** / 3236 / 6004 |
| **T3 コンテナ起動** | **10794** / 10662 / 14473 | **26905** / 26303 / 32555 |
| **T4 ジョブ全体** | **32143** / 26477 / 37536 | **49569** / 45992 / 53783 |
| S 復元後サイズ (B) | 619,411,968 | 604,794,917 |
| S キャッシュサイズ (B) | 151,791,389 | 155,840,397 |

- T1+T2 の中央値差（A − B）: **−11,489 ms**
- T3 の中央値差（B − A）: **+16,111 ms**
- T4 の中央値差（B − A）: **+17,426 ms**

## 判定基準と結論

| 結果 | 判断 | 実測 |
|---|---|---|
| T1+T2 の中央値が方式Aより 20 秒以上短縮、かつ T3 に悪化なし | 採用を検討 | 該当せず |
| 短縮が 10 秒未満 | 見送り | 該当せず（11.5 秒短縮） |
| T2 は消えたが T1 が増えて相殺 | 見送り | 該当せず（T1 はほぼ同じ） |
| **T3 が悪化** | **見送り** | **該当（+16.1 秒）** |
| Phase 1 が通らない | 方式B自体を中止 | 該当せず（Phase 1 は success） |

### 結論: 方式Bは見送り

イメージ投入（T2）は 12.0 秒 → 1.4 秒に減り、`docker load` の逐次展開が
`tar -x` より遅いという仮説自体は裏付けられた。T1 は両者ほぼ同じ。

しかしそれを上回る劣化が2か所で出た。

- **T3（コンテナ起動）が 10.8 秒 → 26.9 秒**。5 回とも一貫している。
- **T0（podman ソケット起動）に 12.1 秒**。方式Aには対応物が無い純増。

結果としてジョブ全体（T4）は方式Bのほうが **17.4 秒遅い**。
イメージ準備を 11.5 秒縮めるために、起動側で 28 秒以上払っている。

T3 悪化の原因は特定していない。方式Bでは `tar -x` で完全に展開済みであり、
「展開の後ろ倒し」では説明がつかない。rootless podman のコンテナ起動コスト
（overlay のマウント、ネットワーク設定、ポートフォワード）が Docker daemon より
重いためと推測されるが、**これは推測であり未検証**。切り分けるなら、
追加ストアを使わず podman に直接 pull させた状態で T3 を測ればよい。

なお当初の計画にあった前提「方式Bでは load 相当の処理が発生しない」は成立しなかった。
actions/cache が必ず tar するため展開は必ず発生し、比較の実体は
「`docker load` 対 `tar -x`」だった（詳細は上記「実行して判明した事実」）。

T0（podman ソケット起動）は方式Bのみが払うコストであり、方式Aには対応物が無い。
headline の比較は T1+T2 で行い、T0 は参考値として別に記録する。
実運用に載せる際の総コストを見るときは T0 も加算すること。

## 既知のリスク

| 項目 | 内容 | 本構成での対処 |
|---|---|---|
| read-only FS の lock file | 追加ストアを read-only FS に置くと lock file 作成で失敗する（containers/storage #1733 ほか） | 書き込み可能な `$HOME/imgstore` に置き、read-only マウントはしない |
| ドライバ不一致 | 主ストアと追加ストアで storage driver が一致しないと参照できない | 生成時 `--storage-driver overlay`、参照時 `driver = "overlay"` を明示 |
| rootless overlay | native overlay が使えず fuse-overlayfs にフォールバックすると遅い | composite action が `graphDriverName` と `mount_program` をログ出力 |
| user systemd | ランナーで `systemctl --user` が使えるか未確認 | 失敗時は `podman system service` に自動フォールバック |
| `XDG_RUNTIME_DIR` 未設定 | ランナーの bash では未設定のことがある | `/run/user/$(id -u)` を試し、不可なら `$HOME/.xdg-runtime` |
| キャッシュサイズ | 展開済みストアは非圧縮 tar より大きい可能性 | S を測定。リポジトリのキャッシュ上限（10 GB）と突き合わせる |
| CI/ローカルの乖離 | CI のみ Podman になる | 環境変数のみで分岐。テストコードには Podman 依存を書かない（`tc.WithProvider` 不使用） |

## 検討したが採らなかった選択肢

| 選択肢 | 除外理由 |
|---|---|
| tar 分割 + `xargs -P` で並列 load | 対象イメージが単一のため並列化の余地がない |
| data-root を tmpfs へ | 2 vCPU / 8 GB。MySQL の書き込みと Go のコンパイルと同居できない |
| `docker load` のバックグラウンド化 | ビルドキャッシュがウォームだと重ねる相手がない |
| skopeo / regctl | `docker-daemon:` transport は daemon の import API 経由で、`docker load` と同一の逐次展開パスを通る |
| レジストリミラー / Docker Hub 認証 / larger runner | 前提として不可 |
| `/var/lib/docker` の直接キャッシュ | 復元時の untar が発生し二重書き込みが解消しない。overlay2 のハードリンク構造と daemon バージョン依存で壊れやすい。方式B不成立時の次善策として保留 |
