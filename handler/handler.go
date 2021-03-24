package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/xerrors"

	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/db"
	. "github.com/ikascrew/ikasbox/handler/internal"
)

func Listen() error {

	err := RegisterStatic()

	err = register()
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	c := config.Get()

	serve := fmt.Sprintf("%s:%d", c.Host, c.Port)
	fmt.Println("ikasbox start[" + serve + "]")

	return http.ListenAndServe(serve, nil)
}

func register() error {

	http.HandleFunc("/", topHandler)
	http.HandleFunc("/content/", contentHandler)
	//http.HandleFunc("/content/upload", contentUploadHandler)
	http.HandleFunc("/content/view/", contentViewHandler)
	http.HandleFunc("/content/media/", contentPlayHandler)

	http.HandleFunc("/group/", groupHandler)
	http.HandleFunc("/group/add", groupAddHandler)
	http.HandleFunc("/group/select", groupSelectHandler)

	http.HandleFunc("/project/", projectListHandler)
	http.HandleFunc("/project/add", projectAddHandler)
	http.HandleFunc("/project/content/list/", projectContentListHandler)

	//TODO 設定から開く
	http.HandleFunc("/thumb/", thumbnailHandler)

	return nil
}

func jsonResponse(w http.ResponseWriter, obj interface{}) error {
	return json.NewEncoder(w).Encode(obj)
}

type Home struct {
	*MenuGroup
}

type MenuGroup struct {
	Selection *db.Group
	List      []*db.Group
}

var gMenuGroup *MenuGroup

func GetMenuGroup() (*MenuGroup, error) {

	if gMenuGroup == nil {
		groups, err := db.SelectGroup()
		if err != nil {
			return nil, xerrors.Errorf("select group: %w", err)
		}
		selection := &db.Group{
			ID:   -1,
			Name: "すべてのコンテンツ",
		}
		gMenuGroup = &MenuGroup{
			Selection: selection,
			List:      groups,
		}
	}

	return gMenuGroup, nil
}

func topHandler(w http.ResponseWriter, r *http.Request) {

	menuGroup, err := GetMenuGroup()
	if err != nil {
		log.Println(err)
	}

	obj := &Home{
		MenuGroup: menuGroup,
	}

	err = Template(w, obj, "top.tmpl")
	if err != nil {
		ErrorPage(w, "Template error", err, 500)

	}

	return
}
