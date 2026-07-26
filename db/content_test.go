package db_test

import (
	"fmt"
	"testing"

	"github.com/ikascrew/ikasbox/db"
)

func newTestContent(t *testing.T, groupId int, name, path string) *db.Content {
	t.Helper()

	c := db.NewContent()
	c.GroupId = groupId
	c.Name = name
	c.Type = "file"
	c.Path = path
	if _, errs := c.Save(); errs != nil {
		t.Fatalf("content save error: %+v", errs)
	}
	return c
}

func TestSelectContentByGroupAndAll(t *testing.T) {
	setupFullDB(t)

	gid1, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}
	gid2, err := db.RegisterGroup("g2", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	for i := 0; i < 3; i++ {
		newTestContent(t, gid1, fmt.Sprintf("c%d", i), fmt.Sprintf("/tmp/c%d.mp4", i))
	}
	newTestContent(t, gid2, "other", "/tmp/other.mp4")

	group1Contents, err := db.SelectContent(gid1)
	if err != nil {
		t.Fatalf("select content by group error: %+v", err)
	}
	if len(group1Contents) != 3 {
		t.Errorf("group1 content count: got %d, want 3", len(group1Contents))
	}

	all, err := db.SelectContent(-1)
	if err != nil {
		t.Fatalf("select all content error: %+v", err)
	}
	if len(all) != 4 {
		t.Errorf("all content count: got %d, want 4", len(all))
	}
}

func TestSelectPagingContent(t *testing.T) {
	setupFullDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}
	otherGid, err := db.RegisterGroup("g2", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	for i := 0; i < 3; i++ {
		newTestContent(t, gid, fmt.Sprintf("c%d", i), fmt.Sprintf("/tmp/c%d.mp4", i))
	}
	newTestContent(t, otherGid, "other", "/tmp/other.mp4")

	pg := &db.Paging{Current: 1, Limit: 2}
	page1, err := db.SelectPagingContent(gid, pg)
	if err != nil {
		t.Fatalf("select paging content (page 1) error: %+v", err)
	}
	if len(page1) != 2 {
		t.Errorf("page 1 length: got %d, want 2", len(page1))
	}
	if pg.Count != 3 {
		t.Errorf("paging count for group: got %d, want 3", pg.Count)
	}

	pg2 := &db.Paging{Current: 2, Limit: 2}
	page2, err := db.SelectPagingContent(gid, pg2)
	if err != nil {
		t.Fatalf("select paging content (page 2) error: %+v", err)
	}
	if len(page2) != 1 {
		t.Errorf("page 2 length: got %d, want 1", len(page2))
	}

	pgAll := &db.Paging{Current: 1, Limit: 10}
	allPage, err := db.SelectPagingContent(-1, pgAll)
	if err != nil {
		t.Fatalf("select paging content (all groups) error: %+v", err)
	}
	if len(allPage) != 4 {
		t.Errorf("all-groups page length: got %d, want 4", len(allPage))
	}
	if pgAll.Count != 4 {
		t.Errorf("all-groups paging count: got %d, want 4", pgAll.Count)
	}
}

func TestDeleteContentRemovesThumbnails(t *testing.T) {
	setupFullDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}
	c := newTestContent(t, gid, "c", "/tmp/c.mp4")

	th := db.ContentThumbnail{ID: c.ID, Seq: 0, Data: []byte("x")}
	if err := th.Insert(); err != nil {
		t.Fatalf("thumbnail insert error: %+v", err)
	}

	if err := db.DeleteContent(c.ID); err != nil {
		t.Fatalf("delete content error: %+v", err)
	}

	if _, err := (db.Content{}).Find(c.ID); err == nil {
		t.Errorf("expected content to be removed")
	}

	ths, err := db.SelectContentThumbnails(c.ID)
	if err != nil {
		t.Fatalf("select thumbnails error: %+v", err)
	}
	if len(ths) != 0 {
		t.Errorf("expected thumbnails removed, got %d", len(ths))
	}
}

func TestDeleteContentNotFound(t *testing.T) {
	setupFullDB(t)

	if err := db.DeleteContent(999); err == nil {
		t.Errorf("expected error deleting missing content")
	}
}
