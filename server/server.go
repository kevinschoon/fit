package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func RunServer(listenPattern, path string) {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", SeriesHandler{dbPath: path, handle: HandleActivities})
	router.Handle("/chart", SeriesHandler{dbPath: path, handle: HandleChart})
	router.HandleFunc("/static/dashboard.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/dashboard.css")
	})
	router.HandleFunc("/static/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/gopher.ico")
	})
	router.HandleFunc("/static/gopher.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/gopher.png")
	})
	log.Printf("Fit server listening @ %s", listenPattern)
	log.Fatal(http.ListenAndServe(listenPattern, handlers.CombinedLoggingHandler(os.Stdout, router)))
}
