package db_test

import (
	"bytes"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/ikascrew/ikasbox/db"
)

//テスト用の一時DBを作成して content_thumbnails テーブルを用意する
func setupThumbnailDB(t *testing.T) {
	t.Helper()

	f, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "thumbnail.db"))
	if err != nil {
		t.Fatalf("initialize db error: %+v", err)
	}
	t.Cleanup(func() { f.Close() })

	_, err = f.Exec(db.CreateContentThumbnailsSQL)
	if err != nil {
		t.Fatalf("create table error: %+v", err)
	}

	db.Use(f)
}

func insertSample(t *testing.T) []byte {
	t.Helper()

	data, err := os.ReadFile("sample.jpg")
	if err != nil {
		t.Fatalf("read file error: %+v", err)
	}

	th := db.ContentThumbnail{}
	th.ID = 1
	th.Seq = 1
	th.Data = data

	err = th.Insert()
	if err != nil {
		t.Fatalf("insert error: %+v", err)
	}

	return data
}

func TestThumbnail(t *testing.T) {

	setupThumbnailDB(t)
	insertSample(t)

	th := db.ContentThumbnail{}
	th.ID = 1
	th.Seq = 1

	err := th.Load()
	if err != nil {
		t.Errorf("load error: %+v", err)
	}
}

func TestSelectContentThumbnails(t *testing.T) {

	setupThumbnailDB(t)
	want := insertSample(t)

	ths, err := db.SelectContentThumbnails(1)
	if err != nil {
		t.Fatalf("select content thumbnails error: %+v", err)
	}

	if len(ths) != 1 {
		t.Fatalf("select content thumbnail count error: got %d", len(ths))
	}

	if !bytes.Equal(ths[0].Data, want) {
		t.Errorf("thumbnail data mismatch: got %d bytes, want %d bytes", len(ths[0].Data), len(want))
	}
}
