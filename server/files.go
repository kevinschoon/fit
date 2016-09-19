package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type StaticHandler struct {
	Path    string   // Path to static assets directory
	Allowed []string // Array of allowed directories
}

func (handler StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, directory := range handler.Allowed {
		if mux.Vars(r)["directory"] == directory {
			if files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", handler.Path, directory)); err == nil {
				for _, file := range files {
					if mux.Vars(r)["file"] == file.Name() && file.Mode().IsRegular() {
						http.ServeFile(w, r, fmt.Sprintf("%s/%s/%s", handler.Path, directory, file.Name()))
						return
					}
				}
			}
		}
	}
	http.NotFound(w, r)
}
