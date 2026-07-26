package db_test

import (
	"testing"

	"github.com/ikascrew/ikasbox/db"
)

func TestRegisterProjectDefaultsResolution(t *testing.T) {
	setupFullDB(t)

	if err := db.RegisterProject("p1", 0, 0); err != nil {
		t.Fatalf("register project error: %+v", err)
	}

	list, err := db.SelectProjectList()
	if err != nil {
		t.Fatalf("select project list error: %+v", err)
	}
	if len(list) != 1 {
		t.Fatalf("project list length: got %d, want 1", len(list))
	}

	p := list[0]
	if p.Name != "p1" {
		t.Errorf("project name: got %q, want p1", p.Name)
	}
	if p.Width != 1280 || p.Height != 720 {
		t.Errorf("default resolution: got %dx%d, want 1280x720", p.Width, p.Height)
	}
}

func TestRegisterProjectExplicitResolution(t *testing.T) {
	setupFullDB(t)

	if err := db.RegisterProject("p1", 640, 480); err != nil {
		t.Fatalf("register project error: %+v", err)
	}

	list, err := db.SelectProjectList()
	if err != nil {
		t.Fatalf("select project list error: %+v", err)
	}
	if len(list) != 1 {
		t.Fatalf("project list length: got %d, want 1", len(list))
	}
	if list[0].Width != 640 || list[0].Height != 480 {
		t.Errorf("resolution: got %dx%d, want 640x480", list[0].Width, list[0].Height)
	}
}

func TestSelectProjectNotFound(t *testing.T) {
	setupFullDB(t)

	if _, err := db.SelectProject(999); err == nil {
		t.Errorf("expected error for missing project")
	}
}

func TestRenameAndDeleteProject(t *testing.T) {
	setupFullDB(t)

	if err := db.RegisterProject("p1", 0, 0); err != nil {
		t.Fatalf("register project error: %+v", err)
	}
	list, err := db.SelectProjectList()
	if err != nil || len(list) != 1 {
		t.Fatalf("select project list error: %+v (len=%d)", err, len(list))
	}
	id := list[0].ID

	if err := db.RenameProject(id, "renamed"); err != nil {
		t.Fatalf("rename project error: %+v", err)
	}

	got, err := db.SelectProject(id)
	if err != nil {
		t.Fatalf("select project error: %+v", err)
	}
	if got.Name != "renamed" {
		t.Errorf("project name after rename: got %q, want renamed", got.Name)
	}

	if err := db.DeleteProject(id); err != nil {
		t.Fatalf("delete project error: %+v", err)
	}

	if _, err := db.SelectProject(id); err == nil {
		t.Errorf("expected error for deleted project")
	}
}

func TestSelectPagingProject(t *testing.T) {
	setupFullDB(t)

	for _, name := range []string{"p1", "p2", "p3"} {
		if err := db.RegisterProject(name, 0, 0); err != nil {
			t.Fatalf("register project error: %+v", err)
		}
	}

	pg := &db.Paging{Current: 1, Limit: 2}
	page1, err := db.SelectPagingProject(pg)
	if err != nil {
		t.Fatalf("select paging project (page 1) error: %+v", err)
	}
	if len(page1) != 2 {
		t.Errorf("page 1 length: got %d, want 2", len(page1))
	}
	if pg.Count != 3 {
		t.Errorf("paging count: got %d, want 3", pg.Count)
	}

	pg2 := &db.Paging{Current: 2, Limit: 2}
	page2, err := db.SelectPagingProject(pg2)
	if err != nil {
		t.Fatalf("select paging project (page 2) error: %+v", err)
	}
	if len(page2) != 1 {
		t.Errorf("page 2 length: got %d, want 1", len(page2))
	}
}
