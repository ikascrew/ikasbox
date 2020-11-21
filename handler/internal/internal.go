package internal

import (
	"fmt"
	"log"
	"net/http"
)

func ErrorPage(w http.ResponseWriter, title string, err error, no int) {

	msg := fmt.Sprintf("%+v", err)
	log.Println(msg)
	dto := struct {
		Title       string
		Description string
		No          int
	}{title, msg, no}

	Template(w, dto, "error.tmpl")
}
