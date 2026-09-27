# API エンドポイントの追加手順

React フロントエンドが呼ぶ JSON API（`/api/v1/...`）は `handler/api/` にある。
1 エンドポイント = 「リクエストを受ける構造体」＋「レスポンスの構造体」＋「`apiMap` への登録」。

## 仕組み

- `api.go` の `init()` にある `apiMap` が、パス（`/api/` を除いた `v1/groups/view` など）から
  `NewParameterFunc`（パラメータ構造体を作る関数）への対応表になっている。
- リクエストボディの JSON はパラメータ構造体に `json.Unmarshal` される。ボディが空なら何もしない（ゼロ値のまま）。
- `Processing() (Return, error)` の戻り値で応答が決まる。
  - `error` を返す → 500 と `{"error": "..."}`
  - `Return.IsSuccess()` が false → 500
  - 成功 → `Return` をそのまま JSON にして 200
- 受け付けるメソッドは POST / PATCH / DELETE のみ（GET は 405）。どれを使うかは慣習で、
  参照系は POST、登録・改名は PATCH、削除は DELETE。サーバー側ではメソッドで振り分けない。

## 1. Go 側

対象に合わせて `group.go` / `project.go` / `content.go` のどれかに書く（どれにも属さないなら新しいファイルを作る）。

```go
type GroupFoo struct {
	GroupId int `json:"groupId"`
}

type GroupFooReturn struct {
	Result string `json:"result"`
	Status
}

func newGroupFoo() Parameter {
	var p GroupFoo
	return &p
}

func (p *GroupFoo) Processing() (Return, error) {
	var ret GroupFooReturn
	// ... db パッケージを呼ぶ。エラーは xerrors.Errorf("...: %w", err) で包んで返す
	ret.success = true
	return &ret, nil
}
```

- [ ] JSON のキーは camelCase（`groupId` など）。フロントからそのまま送れる形にする
- [ ] `Return` 構造体に `Status` を埋め込み、成功時に `ret.success = true` を立てる（忘れると必ず 500 になる）
- [ ] `api.go` の `init()` に `apiMap["v1/groups/foo"] = newGroupFoo` を追加する
- [ ] 時間のかかる処理（取り込み・サムネイル生成）は `GroupRegister` のように goroutine に投げてすぐ返す

### 決まりごと

- **`groupId == -1` は「すべてのグループ」**。コンテンツ一覧・詳細系はこの値を特別扱いする
  （`ContentView.Processing`、`db.SelectContent` / `db.SelectPagingContent`）。
  グループ ID を受け取る新しいエンドポイントでも、一覧系ならこの規約に合わせる。
- ページングは `db.Paging` をリクエストで受け取り、件数を埋めてレスポンスに返す（`GroupView` 参照）。

### ホスト側から足す場合（`api.AddEndpoint`）

ikasbox をライブラリとして同居起動するホスト（`ika-server -ikasbox`）は、
ikasbox に持ち込めない機能を `api.AddEndpoint(path, factory)` で登録する
（例: `v1/server/status` / `v1/server/create`）。

- [ ] `Start`（`handler.Listen`）より **前** に登録する。`apiMap` は排他制御していない
- [ ] ikasbox 本体の `apiMap` にはこれらのパスを入れない。フロントは `v1/server/status` が応答するかどうかで
  「Server」ボタンを出し分けている（`Projects.jsx`）

## 2. フロントエンド側

`src/API.js` は `fetch` の薄いラッパーで、`API.post` / `API.patch` / `API.delete` が `{ data }` を返す
（axios と同じ形）。失敗時（4xx/5xx・通信エラー）は reject するので、`.then()` は成功時にしか走らない。

```js
API.post("/api/v1/groups/foo", { groupId: 1 }).then((res) => {
  this.setState({ result: res.data.result });
});
```

- [ ] URL は `/api/` から書く（開発時は Vite が `/api` を :5555 へ転送する）
- [ ] 失敗時の表示が要るなら `.catch()` を付ける

## 3. テスト

`handler/api/` のテストは実際の SQLite（一時ディレクトリ）に対して HTTP を通して呼ぶ。

- [ ] `setupDB(t)`（`testhelpers_test.go`）で DB を用意し、`serve(t, method, path, body)`（`api_test.go`）で呼ぶ
- [ ] 正常系に加え、`Processing` がエラーを返すケースで 500 と `{"error": ...}` になることを確かめる
- [ ] 例: `group_test.go` の `TestGroupRegisterViewRenameDelete`
- [ ] `go test ./handler/api/`
