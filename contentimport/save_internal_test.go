package contentimport

import (
	"database/sql"
	"image"
	"path/filepath"
	"testing"
	"time"

	"github.com/ikascrew/ikasbox/db"
)

// saveContentWithThumbnails が途中で失敗した場合に、
// コンテンツ行がロールバックされて残らないこと(原子性)を検証する
func TestSaveContentWithThumbnailsRollback(t *testing.T) {

	f, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("db open error: %+v", err)
	}
	t.Cleanup(func() { f.Close() })

	db.Use(f)
	if err := db.CreateTables(); err != nil {
		t.Fatalf("create tables error: %+v", err)
	}

	// 新規コンテンツに採番される id=1 のサムネイル行を先に作っておき、
	// トランザクション内のサムネイル INSERT を主キー衝突で失敗させる
	pre := db.ContentThumbnail{ID: 1, Seq: 0, Data: []byte("x")}
	if err := pre.Insert(); err != nil {
		t.Fatalf("pre-insert thumbnail error: %+v", err)
	}

	images := []image.Image{image.NewRGBA(image.Rect(0, 0, 8, 8))}

	now := time.Now()
	c := db.Content{
		GroupId:   1,
		Name:      "rollback",
		Type:      "img",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := saveContentWithThumbnails(&c, images); err == nil {
		t.Fatal("want error from thumbnail insert conflict, got nil")
	}

	if c.ID != 0 {
		t.Errorf("c.ID = %d, want 0 after rollback", c.ID)
	}

	var n int
	if err := f.QueryRow("select count(*) from contents").Scan(&n); err != nil {
		t.Fatalf("count contents error: %+v", err)
	}
	if n != 0 {
		t.Errorf("contents rows = %d, want 0 (content insert must be rolled back)", n)
	}

	if err := f.QueryRow("select count(*) from content_thumbnails").Scan(&n); err != nil {
		t.Fatalf("count thumbnails error: %+v", err)
	}
	if n != 1 {
		t.Errorf("content_thumbnails rows = %d, want 1 (only the pre-inserted row)", n)
	}
}
