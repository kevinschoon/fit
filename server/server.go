package server

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kevinschoon/gofit/database"
	"log"
	"net/http"
	"os"
)

func RunServer(db *database.DB, pattern, path string) {
	templates := []string{
		fmt.Sprintf("%s/html/base.html", path),
		fmt.Sprintf("%s/html/panel.html", path),
		fmt.Sprintf("%s/html/explore.html", path),
		fmt.Sprintf("%s/html/browse.html", path),
	}
	router := mux.NewRouter().StrictSlash(true)
	handler := Handler{db: db, templates: templates}
	router.Handle("/", ErrorHandler(handler.Home))
	router.Handle("/{series}", ErrorHandler(handler.Home))
	router.Handle("/{series}/chart", ErrorHandler(handler.Home))
	router.Handle("/static/{directory}/{file}", StaticHandler{
		Path: path,
		Allowed: map[string][]string{
			"css": []string{
				"app.css",
			},
			"js": []string{
				"app.js",
			},
		},
	})
	log.Printf("Fit server listening @ %s", pattern)
	log.Fatal(http.ListenAndServe(pattern, handlers.CombinedLoggingHandler(os.Stdout, router)))
}
