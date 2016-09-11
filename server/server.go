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
	staticDir       string = "www"
	baseTmpl        string = staticDir + "/base.html"
	chartTmpl       string = staticDir + "/chart.html"
	dataTmpl        string = staticDir + "/data.html"
	collectionsTmpl string = staticDir + "/collections.html"
)

type TemplateData struct {
	Collection  *models.Collection
	Collections []string
	URLBuilder  *URLBuilder
}

type TemplateOptions struct {
	Collection  *models.Collection
	Collections []string
}

// LoadTemplates loads HTML template files
func LoadTemplate(r *http.Request, options *TemplateOptions) (*template.Template, *TemplateData, error) {
	data := &TemplateData{
		URLBuilder: &URLBuilder{
			URL: r.URL,
		},
	}
	switch {
	case options.Collection != nil:
		data.Collection = options.Collection
		data.URLBuilder.Collection = options.Collection.Name
	case len(options.Collections) > 0:
		data.Collections = options.Collections
	}
	template, err := template.ParseFiles(baseTmpl, chartTmpl, dataTmpl, collectionsTmpl)
	if err != nil {
		return nil, nil, err
	}
	return template, data, nil
}

func RunServer(db *database.DB, pattern string) {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", CollectionsHandler{db: db, handle: Collections})
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
