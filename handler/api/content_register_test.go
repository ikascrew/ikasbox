package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

func TestContentSpec(t *testing.T) {
	setupDB(t)

	w := serve(t, http.MethodPost, "/api/v1/contents/spec", "{}")
	if w.Code != http.StatusOK {
		t.Fatalf("spec status: got %d [%s]", w.Code, w.Body.String())
	}

	var res struct {
		Types []struct {
			Type       string `json:"type"`
			Generative bool   `json:"generative"`
			Fields     []struct {
				Name     string `json:"name"`
				Type     string `json:"type"`
				Label    string `json:"label"`
				Required bool   `json:"required"`
			} `json:"fields"`
		} `json:"types"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("response is not json: %+v [%s]", err, w.Body.String())
	}

	if len(res.Types) == 0 {
		t.Fatal("no types returned")
	}

	// 正語彙が揃っていること、生成型の判定が入っていること
	byType := map[string]bool{}
	for _, e := range res.Types {
		byType[e.Type] = e.Generative
		if len(e.Fields) == 0 {
			t.Errorf("type %q has no fields", e.Type)
		}
		for _, f := range e.Fields {
			if f.Name == "" || f.Type == "" {
				t.Errorf("type %q has incomplete field: %+v", e.Type, f)
			}
		}
	}

	for _, want := range []string{"file", "img", "cd", "terminal"} {
		if _, ok := byType[want]; !ok {
			t.Errorf("missing type %q in spec", want)
		}
	}
	if !byType["cd"] || !byType["terminal"] {
		t.Errorf("cd/terminal should be generative: %+v", byType)
	}
	if byType["file"] || byType["img"] {
		t.Errorf("file/img should not be generative: %+v", byType)
	}
}

// registerGroup はコンテンツ登録先のグループを1つ作って ID を返す
func registerGroup(t *testing.T, name string) int {
	t.Helper()

	w := serve(t, http.MethodPost, "/api/v1/groups/register", `{"name":"`+name+`","path":""}`)
	if w.Code != http.StatusOK {
		t.Fatalf("group register status: got %d [%s]", w.Code, w.Body.String())
	}

	var reg struct {
		GroupId int `json:"groupId"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &reg); err != nil {
		t.Fatalf("group register response is not json: %+v", err)
	}
	if reg.GroupId <= 0 {
		t.Fatalf("groupId not returned: %+v", reg)
	}
	return reg.GroupId
}

func TestContentRegisterGenerated(t *testing.T) {
	setupDB(t)

	gid := registerGroup(t, "g-register")

	body := `{"groupId":` + strconv.Itoa(gid) + `,"name":"NewYear","type":"cd",` +
		`"params":"{\"target\":\"2027-01-01T00:00:00+09:00\",\"text\":\"HNY\"}"}`

	w := serve(t, http.MethodPost, "/api/v1/contents/register", body)
	if w.Code != http.StatusOK {
		t.Fatalf("register status: got %d [%s]", w.Code, w.Body.String())
	}

	var res struct {
		Content struct {
			ID      int    `json:"id"`
			GroupId int    `json:"group_id"`
			Name    string `json:"name"`
			Type    string `json:"type"`
			Params  string `json:"params"`
			Width   int    `json:"width"`
			Height  int    `json:"height"`
		} `json:"content"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("response is not json: %+v [%s]", err, w.Body.String())
	}

	if res.Content.ID <= 0 {
		t.Errorf("content id not returned: %+v", res.Content)
	}
	if res.Content.Type != "cd" {
		t.Errorf("type = %q, want %q", res.Content.Type, "cd")
	}
	if res.Content.Params == "" {
		t.Errorf("params should be stored as-is")
	}
	if res.Content.Width <= 0 || res.Content.Height <= 0 {
		t.Errorf("dimensions: got %dx%d", res.Content.Width, res.Content.Height)
	}
}

// 旧語彙("countdown")でも Normalize されて登録できる
func TestContentRegisterNormalizesLegacyType(t *testing.T) {
	setupDB(t)

	gid := registerGroup(t, "g-legacy")

	body := `{"groupId":` + strconv.Itoa(gid) + `,"name":"Legacy","type":"countdown","params":"{\"text\":\"done\"}"}`

	w := serve(t, http.MethodPost, "/api/v1/contents/register", body)
	if w.Code != http.StatusOK {
		t.Fatalf("register status: got %d [%s]", w.Code, w.Body.String())
	}

	var res struct {
		Content struct {
			Type string `json:"type"`
		} `json:"content"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("response is not json: %+v", err)
	}
	if res.Content.Type != "cd" {
		t.Errorf("type = %q, want normalized %q", res.Content.Type, "cd")
	}
}

func TestContentRegisterInvalidParams(t *testing.T) {
	setupDB(t)

	gid := registerGroup(t, "g-invalid")

	// target が日時として解釈できない -> プラグイン生成に失敗する
	body := `{"groupId":` + strconv.Itoa(gid) + `,"name":"Broken","type":"cd","params":"{\"target\":\"not-a-date\"}"}`

	w := serve(t, http.MethodPost, "/api/v1/contents/register", body)
	if w.Code == http.StatusOK {
		t.Fatalf("expected failure for invalid params, got 200 [%s]", w.Body.String())
	}
	if errorBody(t, w) == "" {
		t.Errorf("error message is empty")
	}
}

func TestContentRegisterValidation(t *testing.T) {
	setupDB(t)

	gid := registerGroup(t, "g-validation")

	cases := []struct {
		name string
		body string
	}{
		{"no groupId", `{"name":"x","type":"cd","params":"{}"}`},
		{"no name", `{"groupId":` + strconv.Itoa(gid) + `,"type":"cd","params":"{}"}`},
		{"unknown group", `{"groupId":9999,"name":"x","type":"cd","params":"{}"}`},
		{"unknown type", `{"groupId":` + strconv.Itoa(gid) + `,"name":"x","type":"no-such","params":"{}"}`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w := serve(t, http.MethodPost, "/api/v1/contents/register", c.body)
			if w.Code == http.StatusOK {
				t.Fatalf("expected failure, got 200 [%s]", w.Body.String())
			}
		})
	}
}
