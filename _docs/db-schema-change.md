# DB スキーマ変更の手順

`db/*_gen.go` は argen（`cmd/argen.exe`、2017 年・modules 以前）の生成物だが、
現在の環境では argen は何も再生成できない。列の追加・変更は **`_gen.go` を手で編集** して行う。
生成コードは機械的なパターンなので、既存の列（実例: `contents.params`）をなぞれば漏れなく書ける。

以下、テーブル `contents` に列 `foo_bar`（Go フィールド `FooBar`）を足す場合を例にする。

## 1. モデル定義（手書きファイル `db/content.go`）

- [ ] `CreateContentsSQL` に列を追加する（新規 `init` 用）
- [ ] `Content` 構造体にフィールドを追加する（`json` タグも付ける）

`string` 列は `NOT NULL DEFAULT ''` にしておく。生成コードは `Scan` で素の `string` に読み込むため、
NULL が入っていると読み出しが失敗する。

## 2. 生成コード（`db/content_gen.go`）

列名が現れる箇所は次の 9 つ。**すべて** に同じ位置（既存の列順）で追加する。

- [ ] `newRelation()` の `r.Select(...)` の列名リスト
- [ ] `Build()` のフィールドコピー（`FooBar: p.FooBar,`）
- [ ] `Save()` の insert 側 `Params(map...)`（`"foo_bar": m.FooBar,`）
- [ ] `Save()` の update 側 `Params(map...)`
- [ ] `Update()` の `if !ar.IsZero(p.FooBar) { m.FooBar = p.FooBar }`
- [ ] `UpdateColumns()` の同じブロック
- [ ] `fieldValueByName()` の `case "foo_bar", "contents.foo_bar": return m.FooBar`
- [ ] `fieldPtrByName()` の `case "foo_bar", "contents.foo_bar": return &m.FooBar`
- [ ] `columnNames()` の列名リスト

実例の位置は `git grep -n -w -i params db/content_gen.go` で一覧できる
（`ar.NewInsert(...)`/`ar.NewUpdate(...)` の `.Params(` 2 行は argen のメソッド名なので列とは無関係）。

注意: `Update`/`UpdateColumns` はゼロ値のフィールドを無視する（argen の仕様）。
空文字や 0 に戻す更新はこの経路ではできないので、必要なら手書きファイル側で直接 SQL を書く。

## 3. 既存 DB へのマイグレーション（`db/db.go` の `migrate()`）

`db.Open` は毎回 `migrate()` を実行する。既存の `ikasbox.db` に列が無ければ追加する処理を足す。

- [ ] `PRAGMA table_info(<table>)` で列の有無を確かめてから `ALTER TABLE ... ADD COLUMN` する（冪等にする）
- [ ] `ALTER TABLE` も `NOT NULL DEFAULT ...` 付きにして、既存行にも値が入るようにする
- [ ] `init` 前（テーブルが無い）の DB では何もしない。テーブル名は `[CONTENTS]` のように大文字で作られているので `lower(name)` で比較する

今の `migrate()` は `contents.params` 専用に書かれている。2 列目以降を足すときは、
列の有無の確認を関数に切り出してから使い回す。

## 4. テスト

- [ ] `db/migrate_internal_test.go` にならい、「古いスキーマの DB に `migrate()` を 2 回かけても列が 1 つだけ追加される」テストを足す
- [ ] 読み書きのテスト（`db/content_test.go` など）で新しい列が保存・読み出しできることを確かめる
- [ ] `go test ./db/ ./handler/api/ ./contentimport/`
