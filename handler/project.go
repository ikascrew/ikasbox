package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ikascrew/ikasbox/db"
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

func projectAddHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	project := db.Project{
		Name:      r.FormValue("name"),
		Width:     1280,
		Height:    720,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, arerr := project.Save(true)
	if arerr != nil {
		log.Println(arerr)
	}

	//選択してあるグループでプロジェクトのコンテンツを登録
	gId, err := strconv.Atoi(r.FormValue("group"))
	if err != nil {
		log.Println(err)
	}

	contentList, err := db.SelectContent(gId)
	if arerr != nil {
		log.Println(arerr)
	}

	for _, elm := range contentList {
		content := db.ProjectContent{
			ProjectID: project.ID,
			ContentID: elm.ID,
			Type:      "file",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, arerr := content.Save(false)
		if err != nil {
			log.Println(arerr)
		}
	}

}

type ProjectResponse struct {
	Project  *db.Project   `json:"project"`
	Contents []*db.Content `json:"contents"`
}

func projectContentListHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.String()
	pathS := strings.Split(path, "/")

	id := pathS[4]

	project, contentList, err := db.SelectProjectContentList(id)
	if err != nil {
	}

	res := ProjectResponse{}
	res.Project = project
	res.Contents = make([]*db.Content, len(contentList))

	for idx, elm := range contentList {
		con, err := elm.Content()
		if err != nil {
		}
		res.Contents[idx] = con
	}

	err = jsonResponse(w, res)
	if err != nil {
		log.Println(err)
	}
}
