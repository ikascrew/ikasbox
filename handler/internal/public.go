package internal

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed _assets/public
var embPublics embed.FS
var publics fs.FS

func init() {
	var err error
	publics, err = fs.Sub(embPublics, "_assets/public")
	if err != nil {
		log.Printf("embed public error:%+v\n", err)
	}
}

func RegisterStatic() error {
	fs := http.FileServer(http.FS(publics))
	http.Handle("/js/", fs)
	http.Handle("/images/", fs)
	http.Handle("/css/", fs)
	return nil
}
