package handler

import (
	"fmt"
	"net/http"

	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/handler/api"
	"github.com/ikascrew/ikasbox/handler/internal"
)

func Listen() error {

	err := register()
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	c := config.Get()

	serve := fmt.Sprintf("%s:%d", c.Host, c.Port)
	fmt.Println("ikasbox start[" + serve + "]")

	err = api.Register("/api/")
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return http.ListenAndServe(serve, nil)
}

func register() error {

	http.HandleFunc("/content/media/", contentPlayHandler)
	http.HandleFunc("/thumb/", thumbnailHandler)

	// External ikascrew tools (server, client) fetch project content lists
	// from this route directly; keep it even though the rest of the legacy
	// HTML UI has been replaced by the React SPA + JSON API.
	http.HandleFunc("/project/content/list/", projectContentListHandler)

	return internal.RegisterSPA()
}
