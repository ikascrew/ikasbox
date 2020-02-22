package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/db"

	"golang.org/x/xerrors"
)

const TemplatePath = "templates/"

func Listen() error {

	register()

	c := config.Get()

	serve := c.Host + ":" + c.Port
	log.Println("ikasbox start[" + serve + "]")

	return http.ListenAndServe(serve, nil)
}

func register() {

	http.HandleFunc("/", topHandler)
	http.HandleFunc("/content/", contentHandler)
	http.HandleFunc("/content/upload", contentUploadHandler)
	//http.HandleFunc("/content/view/",contentViewHandler)
	http.HandleFunc("/group/", groupHandler)
	http.HandleFunc("/group/add", groupAddHandler)
	http.HandleFunc("/group/select", groupSelectHandler)

	http.HandleFunc("/project/", projectListHandler)
	http.HandleFunc("/project/add", projectAddHandler)
	http.HandleFunc("/project/content/list/", projectContentListHandler)

	//TODO 設定から開く

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))

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

	err = layoutWriter(w, obj, TemplatePath+"top.tmpl")
	if err != nil {
		log.Println(err)
	}
}

func layoutWriter(w http.ResponseWriter, o interface{}, tmpl ...string) error {

	args := make([]string, len(tmpl)+1)
	for idx, elm := range tmpl {
		args[idx] = elm
	}
	args[len(tmpl)] = TemplatePath + "layout.tmpl"
	t, err := template.New("layout").ParseFiles(args...)
	if err != nil {
		return err
	}

	return t.Execute(w, o)
}
