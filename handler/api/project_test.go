package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ikascrew/ikasbox/db"
)

func TestProjectRegisterViewRenameDelete(t *testing.T) {
	setupDB(t)

	w := serve(t, http.MethodPost, "/api/v1/projects/register", `{"name":"p1","width":0,"height":0}`)
	if w.Code != http.StatusOK {
		t.Fatalf("register status: got %d [%s]", w.Code, w.Body.String())
	}

	w = serve(t, http.MethodPost, "/api/v1/projects/view", `{"paging":{"current":1,"limit":10}}`)
	if w.Code != http.StatusOK {
		t.Fatalf("view status: got %d [%s]", w.Code, w.Body.String())
	}
	var view struct {
		Projects []struct {
			ID     int    `json:"id"`
			Name   string `json:"name"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"projects"`
		Paging struct {
			Count int `json:"count"`
		} `json:"paging"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &view); err != nil {
		t.Fatalf("response is not json: %+v", err)
	}
	if view.Paging.Count != 1 || len(view.Projects) != 1 {
		t.Fatalf("view mismatch: %+v", view)
	}
	if view.Projects[0].Width != 1280 || view.Projects[0].Height != 720 {
		t.Errorf("default resolution: got %dx%d", view.Projects[0].Width, view.Projects[0].Height)
	}
	pid := view.Projects[0].ID

	body := fmt.Sprintf(`{"projectId":%d,"name":"p1-renamed"}`, pid)
	w = serve(t, http.MethodPatch, "/api/v1/projects/rename", body)
	if w.Code != http.StatusOK {
		t.Fatalf("rename status: got %d [%s]", w.Code, w.Body.String())
	}
	p, err := db.SelectProject(pid)
	if err != nil {
		t.Fatalf("select project error: %+v", err)
	}
	if p.Name != "p1-renamed" {
		t.Errorf("project not renamed: got %q", p.Name)
	}

	body = fmt.Sprintf(`{"projectId":%d}`, pid)
	w = serve(t, http.MethodDelete, "/api/v1/projects/delete", body)
	if w.Code != http.StatusOK {
		t.Fatalf("delete status: got %d [%s]", w.Code, w.Body.String())
	}
	if _, err := db.SelectProject(pid); err == nil {
		t.Errorf("project should be deleted")
	}
}

func TestProjectGroupAddRemoveAndView(t *testing.T) {
	setupDB(t)

	if err := db.RegisterProject("p1", 0, 0); err != nil {
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

	body := fmt.Sprintf(`{"projectId":%d,"groupId":%d}`, pid, gid1)
	w := serve(t, http.MethodPatch, "/api/v1/projects/group/add", body)
	if w.Code != http.StatusOK {
		t.Fatalf("group add status: got %d [%s]", w.Code, w.Body.String())
	}

	w = serve(t, http.MethodPost, "/api/v1/projects/group", fmt.Sprintf(`{"projectId":%d}`, pid))
	if w.Code != http.StatusOK {
		t.Fatalf("group view status: got %d [%s]", w.Code, w.Body.String())
	}
	var gv struct {
		Groups    []struct{ ID int } `json:"groups"`
		AllGroups []struct{ ID int } `json:"allGroups"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &gv); err != nil {
		t.Fatalf("response is not json: %+v", err)
	}
	if len(gv.Groups) != 1 || gv.Groups[0].ID != gid1 {
		t.Errorf("project groups mismatch: %+v", gv.Groups)
	}
	if len(gv.AllGroups) != 2 {
		t.Errorf("all groups count: got %d, want 2", len(gv.AllGroups))
	}
	_ = gid2

	body = fmt.Sprintf(`{"projectId":%d,"groupId":%d}`, pid, gid1)
	w = serve(t, http.MethodDelete, "/api/v1/projects/group/remove", body)
	if w.Code != http.StatusOK {
		t.Fatalf("group remove status: got %d [%s]", w.Code, w.Body.String())
	}

	groups, err := db.SelectProjectGroupList(pid)
	if err != nil {
		t.Fatalf("select project group list error: %+v", err)
	}
	if len(groups) != 0 {
		t.Errorf("groups after remove: got %d, want 0", len(groups))
	}
}

func TestProjectContentsView(t *testing.T) {
	setupDB(t)

	if err := db.RegisterProject("p1", 0, 0); err != nil {
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

	c := db.NewContent()
	c.GroupId = gid
	c.Name = "c1"
	c.Type = "file"
	if _, errs := c.Save(); errs != nil {
		t.Fatalf("save content error: %+v", errs)
	}

	w := serve(t, http.MethodPost, "/api/v1/projects/contents", fmt.Sprintf(`{"projectId":%d}`, pid))
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d [%s]", w.Code, w.Body.String())
	}
	var res struct {
		Project struct {
			ID int `json:"id"`
		} `json:"project"`
		Contents []struct {
			Name string `json:"name"`
		} `json:"contents"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("response is not json: %+v", err)
	}
	if res.Project.ID != pid {
		t.Errorf("project id: got %d, want %d", res.Project.ID, pid)
	}
	if len(res.Contents) != 1 || res.Contents[0].Name != "c1" {
		t.Errorf("contents mismatch: %+v", res.Contents)
	}
}

func TestProjectNotFoundErrors(t *testing.T) {
	setupDB(t)

	w := serve(t, http.MethodPost, "/api/v1/projects/group", `{"projectId":9999}`)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("missing project group view: got %d, want %d", w.Code, http.StatusInternalServerError)
	}

	w = serve(t, http.MethodDelete, "/api/v1/projects/delete", `{"projectId":9999}`)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("delete missing project: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
