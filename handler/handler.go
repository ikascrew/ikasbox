package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"unsafe"

	"github.com/ikascrew/ikasbox/db"
)

const TemplatePath = "templates/"

func Listen() error {

	register()

	//Webを起動
	fmt.Println("Start WebServer")

	//TODO 設定から開く

	return http.ListenAndServe("localhost:5555", nil)
}

func register() {

	http.HandleFunc("/", topHandler)
	http.HandleFunc("/content/", contentHandler)
	http.HandleFunc("/content/upload", contentUploadHandler)
	//http.HandleFunc("/content/view/",contentViewHandler)
	http.HandleFunc("/group/", groupHandler)
	http.HandleFunc("/group/add", groupAddHandler)
	http.HandleFunc("/group/select", groupSelectHandler)

	//TODO 設定から開く

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public/"))))

}

func jsonResponse(w http.ResponseWriter, obj interface{}) error {
	//Header
	return json.NewEncoder(w).Encode(obj)
}

func bstring(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func sbytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

type Home struct {
	*MenuGroup
}

type MenuGroup struct {
	Selection *db.Group
	List      []*db.Group
}

var gMenuGroup *MenuGroup

func GetMenuGroup() *MenuGroup {
	if gMenuGroup == nil {
		groups, err := db.SelectGroup()
		if err != nil {
			panic(err)
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
	return gMenuGroup
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	menuGroup := GetMenuGroup()
	obj := &Home{
		MenuGroup: menuGroup,
	}
	err := layoutWriter(w, obj, TemplatePath+"top.tmpl")
	if err != nil {
		panic(err)
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
