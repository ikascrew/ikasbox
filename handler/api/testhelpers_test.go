package api

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/ikascrew/ikasbox/db"
)

// setupDB は全テーブルを備えた一時 SQLite DB を用意し、パッケージの
// グローバル接続(db.Use)をそれに差し替える。エンドポイントのテストは
// この接続越しに実データを読み書きするため、テスト間の状態は
// t.TempDir() で分離される(db パッケージの setupFullDB と同じ方式)
func setupDB(t *testing.T) {
	t.Helper()

	f, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatalf("initialize db error: %+v", err)
	}
	t.Cleanup(func() { f.Close() })

	db.Use(f)

	if err := db.CreateTables(); err != nil {
		t.Fatalf("create tables error: %+v", err)
	}
}
