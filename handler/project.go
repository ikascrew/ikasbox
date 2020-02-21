package handler

import (
	"github.com/ikascrew/ikasbox/db"
	"net/http"
)

func projectListHandler(w http.ResponseWriter, r *http.Request) {

	list, err := db.SelectProject()
	if err != nil {
	}

	menuGroup, err := GetMenuGroup()
	if err != nil {
	}

	dto := struct {
		ProjectList []*db.Project
		MenuGroup   *MenuGroup
	}{list, menuGroup}

	layoutWriter(w, dto, TemplatePath+"project_list.tmpl")
	if err != nil {
	}
}
