package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ikascrew/ikasbox/db"
)

// ProjectResponse is consumed directly by external ikascrew tools (server's
// config.load, client's tool.getContentList import this type from this
// package), so its shape and the /project/content/list/{id} route below must
// stay stable even though the rest of the legacy HTML UI is gone.
type ProjectResponse struct {
	Project  *db.Project   `json:"project"`
	Contents []*db.Content `json:"contents"`
}

func projectContentListHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.String()
	pathS := strings.Split(path, "/")

	if len(pathS) < 5 {
		errorResponse(w, "url error", fmt.Errorf("path[%s]", path), 400)
		return
	}

	id, err := strconv.Atoi(pathS[4])
	if err != nil {
		errorResponse(w, "url error", err, 400)
		return
	}

	project := db.Project{}
	p, err := project.Find(id)
	if err != nil {
		errorResponse(w, "project error", err, 500)
		return
	}

	contentList, err := db.SelectProjectContentList(id)
	if err != nil {
		errorResponse(w, "select project content list error", err, 500)
		return
	}

	res := ProjectResponse{
		Project:  p,
		Contents: contentList,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		errorResponse(w, "json response error", err, 500)
		return
	}
}
