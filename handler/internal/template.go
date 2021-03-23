package internal

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"golang.org/x/xerrors"
)

//go:embed _assets/templates
var embTmpls embed.FS
var tmpls fs.FS

func init() {
	var err error
	tmpls, err = fs.Sub(embTmpls, "_assets/templates")
	if err != nil {
		log.Printf("embed template error:%+v\n", err)
	}
}

func Template(w http.ResponseWriter, o interface{}, tmpl ...string) error {

	ts := append(tmpl, "layout.tmpl")

	t, err := template.New("layout").ParseFS(tmpls, ts...)
	if err != nil {
		return xerrors.Errorf("Template ParseFiles : %w", err)
	}

	return t.Execute(w, o)
}

func ErrorPage(w http.ResponseWriter, title string, e error, no int) {

	msg := fmt.Sprintf("%+v", e)
	log.Println(msg)
	dto := struct {
		Title       string
		Description string
		No          int
	}{title, msg, no}

	err := Template(w, dto, "error.tmpl")
	if err != nil {
		log.Println("ErrorPage() Error:", err)
	}
}
