package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ikascrew/ikasbox/db"
)

func TestGroupRegisterViewRenameDelete(t *testing.T) {
	setupDB(t)

	// register
	w := serve(t, http.MethodPost, "/api/v1/groups/register", `{"name":"g1","path":""}`)
	if w.Code != http.StatusOK {
		t.Fatalf("register status: got %d [%s]", w.Code, w.Body.String())
	}
	var reg struct {
		GroupId int `json:"groupId"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &reg); err != nil {
		t.Fatalf("response is not json: %+v", err)
	}
	if reg.GroupId <= 0 {
		t.Fatalf("groupId not returned: %+v", reg)
	}

	// view (paging)
	w = serve(t, http.MethodPost, "/api/v1/groups/view", `{"paging":{"current":1,"limit":10}}`)
	if w.Code != http.StatusOK {
		t.Fatalf("view status: got %d [%s]", w.Code, w.Body.String())
	}
	var view struct {
		Groups []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"groups"`
		Paging struct {
			Count int `json:"count"`
		} `json:"paging"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &view); err != nil {
		t.Fatalf("response is not json: %+v", err)
	}
	if view.Paging.Count != 1 || len(view.Groups) != 1 {
		t.Fatalf("view mismatch: %+v", view)
	}
	if view.Groups[0].Name != "g1" {
		t.Errorf("group name: got %q, want g1", view.Groups[0].Name)
	}

	// rename
	body := fmt.Sprintf(`{"groupId":%d,"name":"g1-renamed"}`, reg.GroupId)
	w = serve(t, http.MethodPatch, "/api/v1/groups/rename", body)
	if w.Code != http.StatusOK {
		t.Fatalf("rename status: got %d [%s]", w.Code, w.Body.String())
	}

	g, err := db.FindGroup(reg.GroupId)
	if err != nil {
		t.Fatalf("find group error: %+v", err)
	}
	if g.Name != "g1-renamed" {
		t.Errorf("group not renamed: got %q", g.Name)
	}

	// delete
	body = fmt.Sprintf(`{"groupId":%d}`, reg.GroupId)
	w = serve(t, http.MethodDelete, "/api/v1/groups/delete", body)
	if w.Code != http.StatusOK {
		t.Fatalf("delete status: got %d [%s]", w.Code, w.Body.String())
	}
	if _, err := db.FindGroup(reg.GroupId); err == nil {
		t.Errorf("group should be deleted")
	}
}

func TestGroupRegisterWithPathImportsInBackground(t *testing.T) {
	setupDB(t)

	dir := t.TempDir()
	data, err := os.ReadFile(filepath.Join("..", "..", "contentimport", "testdata", "sample.jpg"))
	if err != nil {
		t.Fatalf("read fixture error: %+v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "a.jpg"), data, 0o644); err != nil {
		t.Fatalf("write fixture error: %+v", err)
	}

	body := fmt.Sprintf(`{"name":"g-import","path":%q}`, dir)
	w := serve(t, http.MethodPost, "/api/v1/groups/register", body)
	if w.Code != http.StatusOK {
		t.Fatalf("register status: got %d [%s]", w.Code, w.Body.String())
	}
	var reg struct {
		GroupId int `json:"groupId"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &reg); err != nil {
		t.Fatalf("response is not json: %+v", err)
	}

	// Import runs in a background goroutine; poll until it has fully landed
	// (content row *and* its 17 thumbnails) before returning, otherwise this
	// test's DB cleanup can race the goroutine's still-in-flight thumbnail
	// inserts and log a harmless but noisy "database is closed" error.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		contents, err := db.SelectContent(reg.GroupId)
		if err != nil {
			t.Fatalf("select content error: %+v", err)
		}
		if len(contents) == 1 {
			ths, err := db.SelectContentThumbnails(contents[0].ID)
			if err != nil {
				t.Fatalf("select thumbnails error: %+v", err)
			}
			if len(ths) == 17 {
				return
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("background import did not register the file in time")
}

func TestGroupCheck(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	dir := t.TempDir()
	existing := filepath.Join(dir, "exists.mp4")
	if err := os.WriteFile(existing, []byte("x"), 0o644); err != nil {
		t.Fatalf("write file error: %+v", err)
	}

	ok := db.NewContent()
	ok.GroupId = gid
	ok.Name = "ok"
	ok.Type = "file"
	ok.Path = existing
	if _, errs := ok.Save(); errs != nil {
		t.Fatalf("save content error: %+v", errs)
	}

	missing := db.NewContent()
	missing.GroupId = gid
	missing.Name = "missing"
	missing.Type = "file"
	missing.Path = filepath.Join(dir, "gone.mp4")
	if _, errs := missing.Save(); errs != nil {
		t.Fatalf("save content error: %+v", errs)
	}

	body := fmt.Sprintf(`{"groupId":%d}`, gid)
	w := serve(t, http.MethodPost, "/api/v1/groups/check", body)
	if w.Code != http.StatusOK {
		t.Fatalf("check status: got %d [%s]", w.Code, w.Body.String())
	}

	var res struct {
		Missing    []struct{ ID int }   `json:"missing"`
		Duplicates [][]struct{ ID int } `json:"duplicates"`
		Total      int                  `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("response is not json: %+v", err)
	}
	if res.Total != 2 {
		t.Errorf("total: got %d, want 2", res.Total)
	}
	if len(res.Missing) != 1 {
		t.Errorf("missing count: got %d, want 1", len(res.Missing))
	}
}

func TestGroupContentsView(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	c := db.NewContent()
	c.GroupId = gid
	c.Name = "c1"
	c.Type = "file"
	if _, errs := c.Save(); errs != nil {
		t.Fatalf("save content error: %+v", errs)
	}

	body := fmt.Sprintf(`{"groupId":%d,"paging":{"current":1,"limit":10}}`, gid)
	w := serve(t, http.MethodPost, "/api/v1/groups/contents", body)
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d [%s]", w.Code, w.Body.String())
	}

	var res struct {
		Group struct {
			Name string `json:"name"`
		} `json:"group"`
		Contents []struct {
			Name string `json:"name"`
		} `json:"contents"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("response is not json: %+v", err)
	}
	if res.Group.Name != "g1" {
		t.Errorf("group name: got %q, want g1", res.Group.Name)
	}
	if len(res.Contents) != 1 || res.Contents[0].Name != "c1" {
		t.Errorf("contents mismatch: %+v", res.Contents)
	}
}

func TestGroupContentsViewAllGroupsSentinel(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}
	c := db.NewContent()
	c.GroupId = gid
	c.Name = "c1"
	c.Type = "file"
	if _, errs := c.Save(); errs != nil {
		t.Fatalf("save content error: %+v", errs)
	}

	w := serve(t, http.MethodPost, "/api/v1/groups/contents", `{"groupId":-1,"paging":{"current":1,"limit":10}}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d [%s]", w.Code, w.Body.String())
	}

	var res struct {
		Group struct {
			ID int `json:"id"`
		} `json:"group"`
		Contents []struct{ ID int } `json:"contents"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("response is not json: %+v", err)
	}
	if res.Group.ID != -1 {
		t.Errorf("group id sentinel: got %d, want -1", res.Group.ID)
	}
	if len(res.Contents) != 1 {
		t.Errorf("contents count: got %d, want 1", len(res.Contents))
	}
}

func TestGroupRegisterAndDeleteNotFoundErrors(t *testing.T) {
	setupDB(t)

	w := serve(t, http.MethodPatch, "/api/v1/groups/rename", `{"groupId":9999,"name":"x"}`)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("rename missing group: got %d, want %d", w.Code, http.StatusInternalServerError)
	}

	w = serve(t, http.MethodDelete, "/api/v1/groups/delete", `{"groupId":9999}`)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("delete missing group: got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
