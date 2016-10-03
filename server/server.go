package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kevinschoon/fit/store"
)

type ErrorHandler func(http.ResponseWriter, *http.Request) error

func (fn ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleError(fn(w, r), w, r)
}

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		switch err.(type) {
		case template.ExecError:
		default:
			switch err {
			case store.ErrNotFound:
				http.NotFound(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

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

func RunServer(db *store.DB, pattern, path, version string, demo bool) {
	templates := []string{
		fmt.Sprintf("%s/html/base.html", path),
		fmt.Sprintf("%s/html/panel.html", path),
		fmt.Sprintf("%s/html/explore.html", path),
		fmt.Sprintf("%s/html/browse.html", path),
	}
	router := mux.NewRouter().StrictSlash(true)
	handler := Handler{db: db, templates: templates, defaults: Response{DemoMode: demo, Version: version}}
	router.Handle("/", ErrorHandler(handler.Home))
	router.Handle("/explore", ErrorHandler(handler.Explore))
	router.Handle("/chart", ErrorHandler(handler.Chart)).Methods("GET")
	if demo {
		router.Handle("/1/dataset", ErrorHandler(handler.DatasetAPI)).Methods("GET")
	} else {
		router.Handle("/1/dataset", ErrorHandler(handler.DatasetAPI)).Methods("GET", "POST", "PUT", "DELETE")
	}
	router.Handle("/static/{directory}/{file}", StaticHandler{
		Path: path,
		Allowed: []string{
			"images",
			"css",
			"js",
			"fonts",
		},
	})
	log.Printf("Fit server listening @ %s", pattern)
	log.Fatal(http.ListenAndServe(pattern, handlers.CombinedLoggingHandler(os.Stdout, router)))
}
