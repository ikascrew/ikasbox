package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//DB不要な範囲(メソッド検証・ルーティング・エラー応答)のテスト

func serve(t *testing.T, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()

	h := &Handle{root: "/api/"}

	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, strings.NewReader(body))
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func errorBody(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()

	var res struct {
		Error string `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("error response is not json: %+v [%s]", err, w.Body.String())
	}
	return res.Error
}

func TestServeHTTPMethodNotAllowed(t *testing.T) {

	w := serve(t, http.MethodGet, "/api/v1/groups/delete", "")

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status code: got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
	if allow := w.Header().Get("Allow"); allow != "POST, PATCH, DELETE" {
		t.Errorf("Allow header: got %q", allow)
	}
	if errorBody(t, w) == "" {
		t.Errorf("error message is empty")
	}
}

func TestServeHTTPNotFound(t *testing.T) {

	w := serve(t, http.MethodPost, "/api/v1/unknown", "{}")

	if w.Code != http.StatusNotFound {
		t.Errorf("status code: got %d, want %d", w.Code, http.StatusNotFound)
	}
	if errorBody(t, w) == "" {
		t.Errorf("error message is empty")
	}
}

func TestServeHTTPBadRequest(t *testing.T) {

	w := serve(t, http.MethodPost, "/api/v1/groups/view", "not json")

	if w.Code != http.StatusBadRequest {
		t.Errorf("status code: got %d, want %d", w.Code, http.StatusBadRequest)
	}
	if errorBody(t, w) == "" {
		t.Errorf("error message is empty")
	}
}
