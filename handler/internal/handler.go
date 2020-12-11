package internal

import (
	"html/template"
	"net/http"

	"golang.org/x/xerrors"
)

const TemplatePath = "assets/templates/"

func Template(w http.ResponseWriter, o interface{}, tmpl ...string) error {

	args := make([]string, len(tmpl)+1)
	for idx, elm := range tmpl {
		args[idx] = TemplatePath + elm
	}

	args[len(tmpl)] = TemplatePath + "layout.tmpl"
	t, err := template.New("layout").ParseFiles(args...)
	if err != nil {
		return xerrors.Errorf("Template ParseFiles : %w", err)
	}

	return t.Execute(w, o)
}
