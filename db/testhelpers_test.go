package db_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/ikascrew/ikasbox/db"
)

// setupFullDB は groups/contents/content_thumbnails/projects/project_groups
// 全テーブルを備えた一時 SQLite DB を用意し、パッケージのグローバル接続
// (db.Use)をそれに差し替える。テスト間の状態は t.TempDir() で分離される
func setupFullDB(t *testing.T) {
	t.Helper()

	f, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "full.db"))
	if err != nil {
		t.Fatalf("initialize db error: %+v", err)
	}
	t.Cleanup(func() { f.Close() })

	db.Use(f)

	if err := db.CreateTables(); err != nil {
		t.Fatalf("create tables error: %+v", err)
	}
}
