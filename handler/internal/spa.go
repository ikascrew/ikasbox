package internal

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

// The React SPA build (see frontend/vite.config.js, outDir points here).
// Run `npm run build` in frontend/ before building/running the Go server.
//
//go:embed all:_assets/spa
var embSpa embed.FS
var spaFS fs.FS

func init() {
	var err error
	spaFS, err = fs.Sub(embSpa, "_assets/spa")
	if err != nil {
		log.Printf("embed spa error:%+v\n", err)
	}
}

// RegisterSPA serves the embedded React build, falling back to index.html
// for any path that isn't a real static asset so client-side routing
// (react-router's BrowserRouter) keeps working on a hard refresh/deep link.
func RegisterSPA() error {
	fileServer := http.FileServer(http.FS(spaFS))
	http.Handle("/", spaHandler(fileServer))
	return nil
}

func spaHandler(fileServer http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		name := strings.TrimPrefix(r.URL.Path, "/")

		if name != "" {
			if f, err := spaFS.Open(name); err == nil {
				f.Close()
				fileServer.ServeHTTP(w, r)
				return
			}
		}

		r2 := new(http.Request)
		*r2 = *r
		r2.URL.Path = "/"
		fileServer.ServeHTTP(w, r2)
	})
}
