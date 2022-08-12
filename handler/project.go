package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ikascrew/ikasbox/db"
	. "github.com/ikascrew/ikasbox/handler/internal"
)

func projectListHandler(w http.ResponseWriter, r *http.Request) {

	// /project/
	// /project/{project id}
	// /project/{project id}/group/{{group id}}

	path := r.URL.String()
	pathS := strings.Split(path, "/")

	menuGroup, err := GetMenuGroup()
	if err != nil {
		ErrorPage(w, "menu group error", err, 500)
		return
	}

	viewId := 0
	addGroup := 0

	if len(pathS) >= 3 {
		idbuf := pathS[2]
		if idbuf != "" {
			viewId, err = strconv.Atoi(idbuf)
			if err != nil {
				ErrorPage(w, "url error", err, 500)
				return
			}
			if len(pathS) >= 5 {
				funcId := pathS[3]
				if funcId == "group" {
					gidBuf := pathS[4]
					addGroup, err = strconv.Atoi(gidBuf)
					if err != nil {
						ErrorPage(w, "url error", err, 500)
						return
					}
				}
			}
		}
	}

	if addGroup != 0 {

		err := db.AddProjectGroup(viewId, addGroup)
		if err != nil {
			ErrorPage(w, "AddProjectGroup() error", err, 500)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/project/%d", viewId), 302)
		return
	}

	if viewId != 0 {

		project, err := db.SelectProject(viewId)
		if err != nil {
			ErrorPage(w, "select project error", err, 500)
			return
		}

		list, err := db.SelectProjectGroupList(viewId)
		if err != nil {
			ErrorPage(w, "select project group list error", err, 500)
			return
		}

		dto := struct {
			Project   *db.Project
			GroupList []*db.Group
			MenuGroup *MenuGroup
		}{project, list, menuGroup}

		err = Template(w, dto, "project_groups.tmpl")
		if err != nil {
			ErrorPage(w, "template error", err, 500)
			return
		}
		return
	}

	list, err := db.SelectProjectList()
	if err != nil {
		ErrorPage(w, "select project error", err, 500)
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

	//a名称を取得してプロジェクトを作成
	r.ParseForm()
	name := r.FormValue("name")
	width := 1280
	height := 720

	project := db.Project{
		Name:      name,
		Width:     width,
		Height:    height,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, arerr := project.Save(true)
	if arerr != nil {
		ErrorPage(w, "project Save() error", arerr, 500)
		return
	}

	obj := struct {
		Success bool
	}{true}

	jsonResponse(w, obj)
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
