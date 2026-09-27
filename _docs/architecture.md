# アーキテクチャ詳細

`AGENTS.md` の Architecture 節の詳細版。パッケージごとの役割と、触るときに知っておくべき経緯をまとめる。

## Go 側

### `cmd/main.go` / パッケージルート

- `cmd/main.go` — CLI の入口。フラグ・引数を解析して `ikasbox.Start()` を呼ぶ。
- `ikasbox.go` — サブコマンド（`start` / `group` / `project` / `init`）を振り分ける。
  `start` は DB を開き、マルチキャストサーバー（`github.com/ikascrew/core/multicast`、サーバー名 "ikasbox"）を
  goroutine で起動してから HTTP サーバーを立てる。
- `group.go` / `project.go`（ルート）— CLI サブコマンドの実装。

### `contentimport/`

ディレクトリ走査・登録・GoCV サムネイル生成のロジック。CLI と Web API の両方から使う。

- CLI（`ikasbox group import`）は対話的な確認とプログレスバーを付けて呼ぶ。
- Web API（`v1/groups/import`）はバックグラウンドの goroutine で走らせてすぐ返す（進捗は返さない）。

関数:

- `ImportDirectory(groupId, path, extensions)` — 一括で取り込む入口。
- `SearchFiles` / `RegisterFile` — ファイル単位で制御したいとき（プログレスバーなど）に組み合わせる部品。
- `RegisterGenerated(groupId, name, type, params)` — 実体ファイルを持たない **生成型コンテンツ**（`cd` / `terminal` など）を登録する。
  `github.com/ikascrew/plugin/video` で実際にプラグインを生成して JSON の `params` を検証し、
  プラグインの描画フレームから 17 枚のサムネイルを作る。

コンテンツの型はプラグインレジストリの正語彙（`file` / `img` / `cd` / `terminal`）を使う。
旧来の `"image"` の行は `video.Normalize` が吸収する。

### `config/`

functional options によるグローバルシングルトン（`config.Set(opts...)` / `config.Get()`）。
既定値はポート 5555、DB `ikasbox.db`。

### `db/`

SQLite 上のデータ層。

- `_gen.go` で終わるファイルは **argen**（`cmd/argen.exe`、ActiveRecord 風 ORM）が生成したもの。
  手書きのロジックは対応する非 `_gen` ファイルに置く（例: `db/group.go` と `db/group_gen.go`）。
- `cmd/argen.exe`（2017 年、modules 以前）はこの環境ではもう再生成できない。スキーマ変更は `_gen` ファイルを手で直す。
  手順は [`db-schema-change.md`](db-schema-change.md)。
- テーブル: `groups`、`contents`（生成型コンテンツ用の `params` TEXT を含む）、`content_thumbnails`、`projects`、`project_groups`。
- `db.Open` は冪等な `migrate()` を実行し、既存 DB に足りない列を `ALTER TABLE` で追加する。
- `db.Transaction(fn)` は関数を begin / commit / rollback で包む。

### `handler/`

3 つのルートと、SPA・API のマウントを持つ。

- `/content/media/{id}` — 動画のバイト列をストリームする（`contentPlayHandler`）。
- `/thumb/{id}[/{seq}]` — サムネイル JPEG を返す（`thumbnailHandler`）。
- `/project/content/list/{id}` — `ProjectResponse{Project, Contents}` の JSON を返す（`projectContentListHandler`、`handler/project.go`）。

#### `/project/content/list/{id}` は外部リポジトリが使っている

このルートは旧 HTML UI の一部ではない。**兄弟リポジトリの `github.com/ikascrew/server` と `github.com/ikascrew/client` が
`handler.ProjectResponse` を直接 import し、この URL をそのまま叩いている**
（`server/config/config.go` の `load()`、`client/tool/project.go` の `getContentList()`）。
`create` 時にローカルのプロジェクト・コンテンツのキャッシュを作るためのもの。
ルートや型を消したり形を変えたりするときは、この 2 リポジトリも同時に直すこと。

それ以外のかつてサーバー側で HTML を描画していた画面（トップ、グループ・プロジェクトの一覧と追加、コンテンツの一覧と詳細）は、
すべて React SPA と `handler/api` に置き換わっている。

### `handler/internal/spa.go`

React の本番ビルド（`_assets/spa`、`npm run build` が出力する）を `go:embed` して `/` で配信する。
実在する静的ファイル以外のパスは `index.html` にフォールバックさせ、react-router のクライアント側ルートが
リロードや直リンクでも動くようにしている。

### `handler/api/`

React フロントエンド向けの JSON API。

- ルーティングは `api.go` の `init()` にあるマップで、パス（例: `v1/groups/view`）から `Parameter` のファクトリへ対応づける。
- 各エンドポイントは `Parameter.Processing() (Return, error)` を実装する構造体で、リクエストの JSON がその構造体に unmarshal される。
- コンテンツの一覧・詳細系は `groupId == -1` を「すべてのグループ」の番兵として使う。
- `api.AddEndpoint(path, factory)` は公開の拡張ポイント。`ika-server -ikasbox`（同居モード）が
  `v1/server/status` / `v1/server/create` を登録するのに使っている。

エンドポイントの追加手順（Go・フロント・テスト）は [`api-endpoint.md`](api-endpoint.md)。

### Vite の開発プロキシ

Vite の開発プロキシは `/content/media` という前方一致だけを転送し、`/content` は転送しない。
`/content` 全体を転送すると、フロント自身の `/contents` SPA ルートまで吸い込んでしまうため。

## React 側（`frontend/src/`）

### 構成

- クラスベースの React 18 コンポーネント + MUI（Material UI）。ルーターは react-router v7
  （統合された `react-router` パッケージ。`react-router-dom` ではない）。
- `API.js` は `fetch` の薄いラッパー（axios 風の `{ data }` を返す）。呼び出しはすべて `/api/v1/...` への
  POST / PATCH / DELETE の JSON。
- Vite でビルドする（`vite.config.js`、入口は `index.html` → `src/App.jsx`）。
- `App.jsx` がアプリをマウントし、`pages/Pages.jsx` がルートを定義する。
- `pages/Layout.jsx`（AppBar + 常設の `Drawer`）と `LayoutDialog.jsx` が共通の枠。
- 機能ごとのページは `pages/Projects/`、`pages/Groups/`、`pages/Contents/`。共通部品は `components/`。

### ファイルの拡張子

JSX を含むファイルは `.jsx`、素のモジュール（`API.js`、`Util.js`、`Paging.js`）は `.js`。
Vite / rollup は JSX を `.jsx` / `.tsx` にしか許さないので、新しいファイルもこれに従う。

### スタイル

- `src/theme.js` が唯一の共有 MUI `createTheme()`。`App.jsx` が全体を `ThemeProvider` + `CssBaseline` で包む。
- 配色・文字・形状はコンポーネントごとに上書きせず、ここで調整する。
- 独自のレイアウト CSS はもう無い（`css/App.css` / `css/Layout.css` は削除済み）。スタイルはテーマと MUI の `sx` プロップで当てる。

### テスト

- テストは対象と同じ場所に置く `*.test.js` / `*.test.jsx`（例: `pages/Paging.test.js`、`components/LoadingButton.test.jsx`）。
- `vitest.config.js` で実行する（jsdom 環境、`src/setupTests.js` が `@testing-library/jest-dom` を読み込む）。
- 純粋なロジック（`Paging.js`、`Util.js`）は素の単体テスト、ライブの API を要しないコンポーネントは
  `@testing-library/react` の描画テストにする。
- マウント時に `API.*` を呼ぶコンポーネントはまだテストされていない。書くなら先に `API.js`
  （または `API.test.js` のように `fetch`）をモックする必要がある。
