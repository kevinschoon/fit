package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kevinschoon/fit/database"
)

func RunServer(db *database.DB, pattern, path, version string, demo bool) {
	templates := []string{
		fmt.Sprintf("%s/html/base.html", path),
		fmt.Sprintf("%s/html/panel.html", path),
		fmt.Sprintf("%s/html/explore.html", path),
		fmt.Sprintf("%s/html/browse.html", path),
	}
	router := mux.NewRouter().StrictSlash(true)
	handler := Handler{db: db, templates: templates, defaults: Response{DemoMode: demo, Version: version}}
	router.Handle("/", ErrorHandler(handler.Home))
	router.Handle("/{series}", ErrorHandler(handler.Home))
	router.Handle("/{series}/chart", ErrorHandler(handler.Chart))
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
