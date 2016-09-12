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
		fmt.Sprintf("%s/html/chart.html", path),
		fmt.Sprintf("%s/html/data.html", path),
		fmt.Sprintf("%s/html/panel.html", path),
		fmt.Sprintf("%s/html/collections.html", path),
	}
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", Handler{db: db, handle: RenderHome, Templates: templates})
	router.Handle("/{collection}", Handler{db: db, handle: RenderHome, Templates: templates})
	router.Handle("/{collection}/chart", Handler{db: db, handle: Chart})
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
