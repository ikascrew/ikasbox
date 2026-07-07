package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"golang.org/x/xerrors"
)

type NewParameterFunc func() Parameter

var apiMap map[string]NewParameterFunc

func init() {
	apiMap = make(map[string]NewParameterFunc)
	apiMap["v1/groups/view"] = newGroupView
	apiMap["v1/groups/register"] = newGroupRegister
	apiMap["v1/groups/rename"] = newGroupRename
	apiMap["v1/groups/check"] = newGroupCheck
	apiMap["v1/groups/delete"] = newGroupDelete
	apiMap["v1/groups/contents"] = newContentView

	apiMap["v1/contents/view"] = newContentFind

	apiMap["v1/projects/view"] = newProjectView
	apiMap["v1/projects/group"] = newProjectGroupView
	apiMap["v1/projects/group/add"] = newProjectGroupAdd
	apiMap["v1/projects/group/remove"] = newProjectGroupRemove
	apiMap["v1/projects/contents"] = newProjectContentView
	apiMap["v1/projects/register"] = newProjectRegister
	apiMap["v1/projects/rename"] = newProjectRename
	apiMap["v1/projects/delete"] = newProjectDelete
}

type Handle struct {
	root string
}

func Register(path string) error {
	var h Handle
	h.root = path

	http.Handle(path, &h)

	return nil
}

//Bodyからデータを抜き出す
func (h *Handle) bind(r *http.Request, p Parameter) error {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return xerrors.Errorf("RequestBody read error(%s): %w", r.URL, err)
	}

	if len(b) == 0 {
		return nil
	}

	err = json.Unmarshal(b, p)
	if err != nil {
		return xerrors.Errorf("json.Unmarshal() error: %w", err)
	}

	return nil
}

func (h *Handle) create(url string) (Parameter, error) {

	np := strings.Replace(url, h.root, "", 1)
	p, ok := apiMap[np]
	if !ok {
		return nil, xerrors.Errorf("not found api[%s]", url)
	}

	return p(), nil
}

func (h *Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	dump(r)

	//フロントが使用するメソッドのみ許可(GETで更新系APIが実行できてしまうのを防ぐ)
	switch r.Method {
	case http.MethodPost, http.MethodPatch, http.MethodDelete:
	default:
		w.Header().Set("Allow", "POST, PATCH, DELETE")
		writeErrorJSON(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed[%s]", r.Method))
		return
	}

	//リクエストからパラメータ引数を指定
	p, err := h.create(r.URL.Path)
	if err != nil {
		writeErrorJSON(w, http.StatusNotFound, err)
		return
	}

	//request body
	err = h.bind(r, p)
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	// schema
	res, err := p.Processing()
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	write(w, res)
}

func write(w http.ResponseWriter, res Return) {

	if res == nil {
		writeErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("Processing() is nil"))
		return
	}

	if !res.IsSuccess() {
		writeErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("processing did not succeed"))
		return
	}

	writeJSON(w, res)
}

func writeJSON(w http.ResponseWriter, x interface{}) {
	//json write
	b, err := json.Marshal(x)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		log.Printf("Write error: %+v", err)
		return
	}
}

type Parameter interface {
	Processing() (Return, error)
}

type Return interface {
	IsSuccess() bool
}

//エラー時はHTTPステータス + {"error": "..."} を返し、フロント側(axios)で検知できるようにする
func writeErrorJSON(w http.ResponseWriter, code int, err error) {

	log.Printf("api error: %+v", err)

	b, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{Error: err.Error()})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(b)
	if err != nil {
		log.Printf("Write error: %+v", err)
	}
}

func dump(r *http.Request) {

	if true {
		return
	}

	dump, _ := httputil.DumpRequest(r, true)
	line := "-------------------------------------------------------"
	log.Printf("%s\n%s\n", line, string(dump))
	log.Println(line)
}
