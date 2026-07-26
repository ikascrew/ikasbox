package db_test

import (
	"testing"

	"github.com/ikascrew/ikasbox/db"
)

func TestProjectGroupAddRemoveAndContentList(t *testing.T) {
	setupFullDB(t)

	if err := db.RegisterProject("proj", 0, 0); err != nil {
		t.Fatalf("register project error: %+v", err)
	}
	projects, err := db.SelectProjectList()
	if err != nil || len(projects) != 1 {
		t.Fatalf("select project list error: %+v (len=%d)", err, len(projects))
	}
	pid := projects[0].ID

	gid1, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}
	gid2, err := db.RegisterGroup("g2", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	newTestContent(t, gid1, "a", "/tmp/a.mp4")
	newTestContent(t, gid2, "b", "/tmp/b.mp4")

	if err := db.AddProjectGroup(pid, gid1); err != nil {
		t.Fatalf("add project group error: %+v", err)
	}
	if err := db.AddProjectGroup(pid, gid2); err != nil {
		t.Fatalf("add project group error: %+v", err)
	}

	groups, err := db.SelectProjectGroupList(pid)
	if err != nil {
		t.Fatalf("select project group list error: %+v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("project group list length: got %d, want 2", len(groups))
	}

	contents, err := db.SelectProjectContentList(pid)
	if err != nil {
		t.Fatalf("select project content list error: %+v", err)
	}
	if len(contents) != 2 {
		t.Errorf("project content list length: got %d, want 2", len(contents))
	}

	if err := db.RemoveProjectGroup(pid, gid1); err != nil {
		t.Fatalf("remove project group error: %+v", err)
	}
	groupsAfterRemove, err := db.SelectProjectGroupList(pid)
	if err != nil {
		t.Fatalf("select project group list error: %+v", err)
	}
	if len(groupsAfterRemove) != 1 {
		t.Fatalf("project group list length after remove: got %d, want 1", len(groupsAfterRemove))
	}
	if groupsAfterRemove[0].ID != gid2 {
		t.Errorf("remaining group: got %d, want %d", groupsAfterRemove[0].ID, gid2)
	}

	if err := db.DeleteProjectGroups(pid); err != nil {
		t.Fatalf("delete project groups error: %+v", err)
	}
	groupsAfterDeleteAll, err := db.SelectProjectGroupList(pid)
	if err != nil {
		t.Fatalf("select project group list error: %+v", err)
	}
	if len(groupsAfterDeleteAll) != 0 {
		t.Errorf("project group list length after delete all: got %d, want 0", len(groupsAfterDeleteAll))
	}
}

func TestDeleteProjectRemovesProjectGroups(t *testing.T) {
	setupFullDB(t)

	if err := db.RegisterProject("proj", 0, 0); err != nil {
		t.Fatalf("register project error: %+v", err)
	}
	projects, err := db.SelectProjectList()
	if err != nil || len(projects) != 1 {
		t.Fatalf("select project list error: %+v (len=%d)", err, len(projects))
	}
	pid := projects[0].ID

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	if err := db.AddProjectGroup(pid, gid); err != nil {
		t.Fatalf("add project group error: %+v", err)
	}

	if err := db.DeleteProject(pid); err != nil {
		t.Fatalf("delete project error: %+v", err)
	}

	// group itself (and its contents) must survive project deletion
	if _, err := db.FindGroup(gid); err != nil {
		t.Errorf("group should survive project delete: %+v", err)
	}
}
