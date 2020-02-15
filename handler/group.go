package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ikascrew/ikasbox/db"
)

func groupHandler(w http.ResponseWriter, r *http.Request) {
}

func groupAddHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	group := db.Group{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := group.Save(false)
	if err != nil {
		panic(err)
	}

	obj := struct {
		Success bool
	}{true}

	jsonResponse(w, obj)
}

func groupSelectHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	groupId := r.FormValue("id")
	intId, err := strconv.Atoi(groupId)
	if err != nil {
		panic(err)
	}

	menuGroup := GetMenuGroup()
	var group *db.Group
	for _, elm := range menuGroup.List {
		if elm.ID == intId {
			group = elm
			break
		}
	}

	menuGroup.Selection = group
	obj := struct {
		Success bool
	}{true}
	jsonResponse(w, obj)
}
