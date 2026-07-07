package handler

import (
	"fmt"
	"net/http"

	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/handler/api"
	. "github.com/ikascrew/ikasbox/handler/internal"
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

	return RegisterSPA()
}
