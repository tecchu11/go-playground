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
| T1 | tar の DL + zstd 伸長 | ストアディレクトリの DL + zstd 伸長 |
| T2 | `docker load` | 追加ストア参照の確認 |
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

## 判定基準

| 結果 | 判断 |
|---|---|
| T1+T2 の中央値が方式Aより 20 秒以上短縮、かつ T3 に悪化なし | 採用を検討 |
| 短縮が 10 秒未満 | 見送り |
| T2 は消えたが T1 が増えて相殺 | 見送り。展開コストがキャッシュ復元側に移動しただけ |
| T3 が悪化 | 展開が起動時に後ろ倒しされている。見送り |
| Phase 1 が通らない | 方式B自体を中止 |

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
