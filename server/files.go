package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type StaticHandler struct {
	Path    string // Path to static assets directory
	Allowed map[string][]string
}

func (handler StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for directory, files := range handler.Allowed {
		if mux.Vars(r)["directory"] == directory {
			for _, file := range files {
				if mux.Vars(r)["file"] == file {
					http.ServeFile(w, r, fmt.Sprintf("%s/%s/%s", handler.Path, directory, file))
					return
				}
			}
		}
	}
	http.NotFound(w, r)
}
