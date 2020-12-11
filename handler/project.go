package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ikascrew/ikasbox/db"
	. "github.com/ikascrew/ikasbox/handler/internal"
)

func projectListHandler(w http.ResponseWriter, r *http.Request) {

	list, err := db.SelectProject()
	if err != nil {
		ErrorPage(w, "select project error", err, 500)
		return
	}

	menuGroup, err := GetMenuGroup()
	if err != nil {
		ErrorPage(w, "menu group error", err, 500)
		return
	}

	dto := struct {
		ProjectList []*db.Project
		MenuGroup   *MenuGroup
	}{list, menuGroup}

	err = Template(w, dto, "project_list.tmpl")
	if err != nil {
		ErrorPage(w, "template error", err, 500)
		return
	}
}

func projectAddHandler(w http.ResponseWriter, r *http.Request) {

}

type ProjectResponse struct {
	Project  *db.Project   `json:"project"`
	Contents []*db.Content `json:"contents"`
}

func projectContentListHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.String()
	pathS := strings.Split(path, "/")

	idbuf := pathS[4]
	id, err := strconv.Atoi(idbuf)
	if err != nil {
		ErrorPage(w, "url error", err, 500)
		return
	}

	project := db.Project{}
	p, err := project.Find(id)
	if err != nil {
		ErrorPage(w, "project error", err, 500)
		return
	}

	res := ProjectResponse{}
	res.Project = p

	contentList, err := db.SelectProjectContentList(id)
	if err != nil {
		ErrorPage(w, "select project content list error", err, 500)
		return
	}

	res.Contents = contentList

	err = jsonResponse(w, res)
	if err != nil {
		ErrorPage(w, "json response error", err, 500)
		return
	}
}
