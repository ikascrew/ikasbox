package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

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

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.String()
	us := strings.Split(url, "/")

	if len(us) < 3 {
		ErrorPage(w, "Thumbnail not found", fmt.Errorf("url[%s]", url), 404)
		return
	}

	idbuf := us[2]
	id, err := strconv.Atoi(idbuf)
	if err != nil {
		ErrorPage(w, "thumbnail id error", err, 400)
		return
	}

	seq := 0

	if len(us) >= 4 {
		seqbuf := us[3]
		seq, err = strconv.Atoi(seqbuf)
		if err != nil {
			ErrorPage(w, "thumbnail seq error", err, 400)
			return
		}
	}

	thumb := db.NewContentThumbnail()
	thumb.ID = id
	thumb.Seq = seq

	err = thumb.Load()
	if err != nil {
		ErrorPage(w, "thumbnail database not found", err, 404)
		return
	}

	if thumb.Data == nil {
		ErrorPage(w, "thumbnail data is nil", fmt.Errorf("database error"), 404)
		return
	}

	_, err = w.Write(thumb.Data)
	if err != nil {
		ErrorPage(w, "thumbnail write error", err, 404)
		return
	}
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
