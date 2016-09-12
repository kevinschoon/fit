package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kevinschoon/gofit/database"
	"log"
	"net/http"
	"os"
)

func RunServer(db *database.DB, pattern string) {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", Handler{db: db, handle: RenderHome})
	router.Handle("/{collection}", Handler{db: db, handle: RenderHome})
	router.Handle("/{collection}/chart", Handler{db: db, handle: Chart})
	router.HandleFunc("/static/dashboard.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/dashboard.css")
	})
	router.HandleFunc("/static/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/gopher.ico")
	})
	router.HandleFunc("/static/gopher.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/gopher.png")
	})
	log.Printf("Fit server listening @ %s", pattern)
	log.Fatal(http.ListenAndServe(pattern, handlers.CombinedLoggingHandler(os.Stdout, router)))
}
