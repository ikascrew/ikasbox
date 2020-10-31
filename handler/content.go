package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/ikascrew/ikasbox/db"
	. "github.com/ikascrew/ikasbox/handler/internal"
)

type Content struct {
	ContentList []*db.Content
	*MenuGroup
}

func contentHandler(w http.ResponseWriter, r *http.Request) {

	menu, err := GetMenuGroup()
	if err != nil {
		log.Println(err)
	}

	gId := menu.Selection.ID

	contentList, err := db.SelectContent(gId)
	if err != nil {
		log.Println(err)
	}

	obj := Content{
		ContentList: contentList,
		MenuGroup:   menu,
	}

	err = Template(w, obj, "content.tmpl")
	if err != nil {
		ErrorPage(w, "Template error", err, 500)
	}
}

func contentViewHandler(w http.ResponseWriter, r *http.Request) {

	idBuf := path.Base(r.URL.Path)
	id, err := strconv.Atoi(idBuf)

	if err != nil {
		ErrorPage(w, "Content ID NotFound", err, 404)
		return
	}

	datum, err := db.Content{}.Find(id)
	if err != nil {
		ErrorPage(w, "Content Select Error", err, 500)
		return
	}
	menu, err := GetMenuGroup()
	if err != nil {
		log.Println(err)
	}

	width := datum.Width
	height := datum.Height
	for width > 720 {
		width = width / 2
		height = height / 2
	}

	dto := struct {
		Content   *db.Content
		MenuGroup *MenuGroup
		Width     int
		Height    int
	}{datum, menu, width, height}

	err = Template(w, dto, "content_view.tmpl")
	if err != nil {
		ErrorPage(w, "Template error", err, 500)
	}
}

func contentPlayHandler(w http.ResponseWriter, r *http.Request) {
	idBuf := path.Base(r.URL.Path)

	id, err := strconv.Atoi(idBuf)
	if err != nil {
		ErrorPage(w, "Countent ID not found", err, 404)
		return
	}

	datum, err := db.Content{}.Find(id)
	if err != nil {
		ErrorPage(w, "Content Select Error", err, 500)
		return
	}

	fs, err := os.Open(datum.Path)
	if err != nil {
		fmt.Println(err)
		ErrorPage(w, "Content Open Error", err, 500)
		return
	}
	defer fs.Close()

	_, err = io.Copy(w, fs)
	if err != nil {
		fmt.Println(err)
		ErrorPage(w, "Media Copy Error", err, 500)
		return
	}

	//w.Write(sbytes(datum.Content))
}

/*
func contentUploadHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	name := r.FormValue("content_name")
	//file,fHead,err := r.FormFile("content_datum")

	fmt.Println("content_name=", name)

	now := time.Now()
	menu, err := GetMenuGroup()
	if err != nil {
		ErrorPage(w, "Menu Group Error", 500)
		return
	}

	menuId := menu.Selection.ID

	content := db.Content{
		Name:      name,
		GroupId:   menuId,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = db.Transaction(func(tx *sql.Tx) error {

		_, arErr := content.Save(false)
		if arErr != nil {
			return fmt.Errorf("Content save error:%s", arErr)
		}

				buf := make([]byte,fHead.Size)
				file.Read(buf)

				datum := db.ContentDatum{
					ID:content.ID,
					Content:bstring(buf),
					CreatedAt:now,
					UpdatedAt:now,
				}

				flg,arErr =datum.Save(false)
			    if arErr != nil {
				    return fmt.Errorf("ContentData save error:%s",arErr)
			    }
				fmt.Println(flg)
		return nil
	})

	if err != nil {
		ErrorPage(w, "Database Error", err, 500)
		return
	}

	menuGroup, err := GetMenuGroup()
	if err != nil {
		ErrorPage(w, "Menu Group Error", err, 500)
		return
	}

	obj := Content{
		MenuGroup: menuGroup,
	}
	err = Template(w, obj, "content.tmpl")
	if err != nil {
		ErrorPage(w, "Layout Error", err, 500)
		return
	}
}
*/
