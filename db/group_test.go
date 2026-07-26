package db_test

import (
	"fmt"
	"testing"

	"github.com/ikascrew/ikasbox/db"
)

func TestRegisterFindRenameGroup(t *testing.T) {
	setupFullDB(t)

	id, err := db.RegisterGroup("g1", "/path/a")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}
	if id <= 0 {
		t.Fatalf("register group id: got %d", id)
	}

	g, err := db.FindGroup(id)
	if err != nil {
		t.Fatalf("find group error: %+v", err)
	}
	if g.Name != "g1" || g.Path != "/path/a" {
		t.Errorf("group mismatch: got name=%q path=%q", g.Name, g.Path)
	}

	if err := db.RenameGroup(id, "g1-renamed"); err != nil {
		t.Fatalf("rename group error: %+v", err)
	}

	renamed, err := db.FindGroup(id)
	if err != nil {
		t.Fatalf("find group after rename error: %+v", err)
	}
	if renamed.Name != "g1-renamed" {
		t.Errorf("renamed group name: got %q, want %q", renamed.Name, "g1-renamed")
	}
	// path is left untouched by rename
	if renamed.Path != "/path/a" {
		t.Errorf("group path changed by rename: got %q", renamed.Path)
	}
}

func TestFindGroupNotFound(t *testing.T) {
	setupFullDB(t)

	if _, err := db.FindGroup(999); err == nil {
		t.Errorf("expected error for missing group")
	}
}

func TestDeleteGroupCascadesContentsAndThumbnails(t *testing.T) {
	setupFullDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	c := db.NewContent()
	c.GroupId = gid
	c.Name = "content"
	c.Type = "file"
	if _, errs := c.Save(); errs != nil {
		t.Fatalf("content save error: %+v", errs)
	}

	th := db.ContentThumbnail{ID: c.ID, Seq: 0, Data: []byte("x")}
	if err := th.Insert(); err != nil {
		t.Fatalf("thumbnail insert error: %+v", err)
	}

	if err := db.DeleteGroup(gid); err != nil {
		t.Fatalf("delete group error: %+v", err)
	}

	if _, err := db.FindGroup(gid); err == nil {
		t.Errorf("expected error finding deleted group")
	}

	if _, err := (db.Content{}).Find(c.ID); err == nil {
		t.Errorf("expected content to be removed by cascading group delete")
	}

	ths, err := db.SelectContentThumbnails(c.ID)
	if err != nil {
		t.Fatalf("select thumbnails error: %+v", err)
	}
	if len(ths) != 0 {
		t.Errorf("expected thumbnails removed by cascading delete, got %d", len(ths))
	}
}

func TestSelectGroupAndPaging(t *testing.T) {
	setupFullDB(t)

	const total = 3
	for i := 0; i < total; i++ {
		if _, err := db.RegisterGroup(fmt.Sprintf("g%d", i), ""); err != nil {
			t.Fatalf("register group[%d] error: %+v", i, err)
		}
	}

	all, err := db.SelectGroup()
	if err != nil {
		t.Fatalf("select group error: %+v", err)
	}
	if len(all) != total {
		t.Fatalf("select group count: got %d, want %d", len(all), total)
	}

	pg := &db.Paging{Current: 1, Limit: 2}
	page1, err := db.SelectPagingGroup(pg)
	if err != nil {
		t.Fatalf("select paging group (page 1) error: %+v", err)
	}
	if len(page1) != 2 {
		t.Errorf("page 1 length: got %d, want 2", len(page1))
	}
	if pg.Count != total {
		t.Errorf("paging count: got %d, want %d", pg.Count, total)
	}

	pg2 := &db.Paging{Current: 2, Limit: 2}
	page2, err := db.SelectPagingGroup(pg2)
	if err != nil {
		t.Fatalf("select paging group (page 2) error: %+v", err)
	}
	if len(page2) != 1 {
		t.Errorf("page 2 length: got %d, want 1", len(page2))
	}
}
