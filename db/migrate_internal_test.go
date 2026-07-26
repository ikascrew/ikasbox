package db

import (
	"database/sql"
	"path/filepath"
	"testing"
)

// migrate() is unexported, so this file lives in package db (not db_test)
// to exercise it directly.

func TestMigrateNoContentsTableIsNoop(t *testing.T) {
	f, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "empty.db"))
	if err != nil {
		t.Fatalf("open db error: %+v", err)
	}
	defer f.Close()

	Use(f)

	if err := migrate(); err != nil {
		t.Fatalf("migrate on a db with no contents table should be a no-op: %+v", err)
	}
}

func TestMigrateAddsParamsColumnAndIsIdempotent(t *testing.T) {
	f, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "legacy.db"))
	if err != nil {
		t.Fatalf("open db error: %+v", err)
	}
	defer f.Close()

	Use(f)

	// simulate a pre-"params" column legacy schema
	if _, err := f.Exec(`CREATE TABLE contents (id INTEGER PRIMARY KEY, group_id INTEGER, name TEXT)`); err != nil {
		t.Fatalf("create legacy table error: %+v", err)
	}

	if err := migrate(); err != nil {
		t.Fatalf("migrate error: %+v", err)
	}

	if !hasParamsColumn(t, f) {
		t.Fatalf("params column was not added by migrate")
	}

	// calling migrate again must not error (idempotent)
	if err := migrate(); err != nil {
		t.Fatalf("second migrate call errored: %+v", err)
	}
}

func TestMigrateRepairsNullParams(t *testing.T) {
	f, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "nullparams.db"))
	if err != nil {
		t.Fatalf("open db error: %+v", err)
	}
	defer f.Close()

	Use(f)

	if _, err := f.Exec(`CREATE TABLE contents (id INTEGER PRIMARY KEY, group_id INTEGER, name TEXT, params TEXT)`); err != nil {
		t.Fatalf("create table error: %+v", err)
	}
	if _, err := f.Exec(`INSERT INTO contents (id, group_id, name, params) VALUES (1, 1, 'x', NULL)`); err != nil {
		t.Fatalf("insert row error: %+v", err)
	}

	if err := migrate(); err != nil {
		t.Fatalf("migrate error: %+v", err)
	}

	var params sql.NullString
	if err := f.QueryRow(`SELECT params FROM contents WHERE id = 1`).Scan(&params); err != nil {
		t.Fatalf("select params error: %+v", err)
	}
	if !params.Valid || params.String != "" {
		t.Errorf("params should be repaired to empty string, got valid=%v value=%q", params.Valid, params.String)
	}
}

func hasParamsColumn(t *testing.T, f *sql.DB) bool {
	t.Helper()

	rows, err := f.Query(`PRAGMA table_info(contents)`)
	if err != nil {
		t.Fatalf("table_info(contents): %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var colName, typ string
		var notnull int
		var dflt interface{}
		var pk int
		if err := rows.Scan(&cid, &colName, &typ, &notnull, &dflt, &pk); err != nil {
			t.Fatalf("scan table_info: %+v", err)
		}
		if colName == "params" {
			return true
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("table_info rows: %+v", err)
	}
	return false
}
