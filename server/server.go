package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kevinschoon/gofit/database"
	"github.com/kevinschoon/gofit/models"
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	staticDir string = "www"
	baseTmpl  string = staticDir + "/base.html"
	chartTmpl string = staticDir + "/chart.html"
	dataTmpl  string = staticDir + "/data.html"
)

type TemplateData struct {
	Collection *models.Collection
	URLBuilder *URLBuilder
}

// LoadTemplates loads HTML template files
func LoadTemplate(r *http.Request, collection *models.Collection) (*template.Template, *TemplateData, error) {
	data := &TemplateData{
		Collection: collection,
		URLBuilder: &URLBuilder{
			URL:        r.URL,
			Collection: collection.Name,
		},
	}
	template, err := template.ParseFiles(baseTmpl, chartTmpl, dataTmpl)
	if err != nil {
		return nil, nil, err
	}
	return template, data, nil
}

func RunServer(db *database.DB, pattern string) {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/{collection}", CollectionHandler{db: db, handle: Collection})
	router.Handle("/{collection}/chart", CollectionHandler{db: db, handle: Chart})
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
