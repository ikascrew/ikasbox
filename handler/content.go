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
)

func errorResponse(w http.ResponseWriter, msg string, err error, code int) {
	log.Printf("%s: %+v", msg, err)
	http.Error(w, msg, code)
}

func contentPlayHandler(w http.ResponseWriter, r *http.Request) {
	idBuf := path.Base(r.URL.Path)

	id, err := strconv.Atoi(idBuf)
	if err != nil {
		errorResponse(w, "Countent ID not found", err, 404)
		return
	}

	datum, err := db.Content{}.Find(id)
	if err != nil {
		errorResponse(w, "Content Select Error", err, 500)
		return
	}

	fs, err := os.Open(datum.Path)
	if err != nil {
		fmt.Println(err)
		errorResponse(w, "Content Open Error", err, 500)
		return
	}
	defer fs.Close()

	_, err = io.Copy(w, fs)
	if err != nil {
		fmt.Println(err)
		errorResponse(w, "Media Copy Error", err, 500)
		return
	}
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.String()
	us := strings.Split(url, "/")

	if len(us) < 3 {
		errorResponse(w, "Thumbnail not found", fmt.Errorf("url[%s]", url), 404)
		return
	}

	idbuf := us[2]
	id, err := strconv.Atoi(idbuf)
	if err != nil {
		errorResponse(w, "thumbnail id error", err, 400)
		return
	}

	seq := 0

	if len(us) >= 4 {
		seqbuf := us[3]
		seq, err = strconv.Atoi(seqbuf)
		if err != nil {
			errorResponse(w, "thumbnail seq error", err, 400)
			return
		}
	}

	thumb := db.NewContentThumbnail()
	thumb.ID = id
	thumb.Seq = seq

	err = thumb.Load()
	if err != nil {
		errorResponse(w, "thumbnail database not found", err, 404)
		return
	}

	if thumb.Data == nil {
		errorResponse(w, "thumbnail data is nil", fmt.Errorf("database error"), 404)
		return
	}

	_, err = w.Write(thumb.Data)
	if err != nil {
		errorResponse(w, "thumbnail write error", err, 404)
		return
	}
}
