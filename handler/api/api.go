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

func (h *Handle) createParameter(r *http.Request) (Parameter, error) {

	url := r.URL

	//リクエストからパラメータ引数を指定
	p, err := h.create(url.Path)
	if err != nil {
		return nil, xerrors.Errorf("Handle.create(%s) error: %w", url, err)
	}

	//Bodyからデータを抜き出す
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, xerrors.Errorf("RequestBody read error: %w", url, err)
	}

	if b == nil || len(b) == 0 {
		return p, nil
	}

	err = json.Unmarshal(b, p)
	if err != nil {
		return nil, xerrors.Errorf("json.Unmarshal() error: %w", err)
	}

	return p, nil
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

	//request body
	p, err := h.createParameter(r)
	if err != nil {
		//Not Found
		writeErrorJSON(w, err)
		return
	}

	// schema
	res := p.Processing()

	write(w, res)
}

func write(w http.ResponseWriter, res Return) {

	if res == nil {
		//Not Found
		writeErrorJSON(w, fmt.Errorf("Processing() is nil"))
		return
	}

	if !res.IsSuccess() {
		writeStatusJSON(w, res.GetStatus())
		return
	}

	writeJSON(w, res)
}

func writeJSON(w http.ResponseWriter, x interface{}) {
	//json write
	b, err := json.Marshal(x)
	if err != nil {
		writeErrorJSON(w, err)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		log.Println("Write error: %+v", err)
		return
	}
}

type Parameter interface {
	Processing() Return
}

type Return interface {
	GetStatus() Status
	IsSuccess() bool
}

func writeStatusJSON(w http.ResponseWriter, status Status) {
}

func writeErrorJSON(w http.ResponseWriter, err error) {

	log.Printf("write error: %+v", err)

	res := struct {
	}{}

	writeJSON(w, res)
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
