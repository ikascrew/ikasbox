package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/ikascrew/ikasbox/db"
)

type Content struct {
	ContentList []*db.Content
	*MenuGroup
}

func contentHandler(w http.ResponseWriter, r *http.Request) {

	menu, err := GetMenuGroup()
	if err != nil {
	}

	gId := menu.Selection.ID

	contentList, err := db.SelectContent(gId)
	if err != nil {
		panic(err)
	}

	obj := Content{
		ContentList: contentList,
		MenuGroup:   menu,
	}
	err = layoutWriter(w, obj, TemplatePath+"content.tmpl")
	if err != nil {
	}
}

/*
func contentViewHandler(w http.ResponseWriter,r *http.Request) {
	idBuf := path.Base(r.URL.Path)
	id,err := strconv.Atoi(idBuf)
	if err != nil {
		panic(err)
	}
	datum,err := db.ContentDatum{}.Find(id)
	w.Write(sbytes(datum.Content))
}
*/

func contentUploadHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	name := r.FormValue("content_name")
	//file,fHead,err := r.FormFile("content_datum")

	fmt.Println("content_name=", name)

	now := time.Now()
	menu, err := GetMenuGroup()
	if err != nil {
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

		/*
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
		*/
		return nil
	})

	if err != nil {

	}

	menuGroup, err := GetMenuGroup()
	if err != nil {
	}

	obj := Content{
		MenuGroup: menuGroup,
	}
	err = layoutWriter(w, obj, TemplatePath+"content.tmpl")
	if err != nil {
		panic(err)
	}
}
