package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kevinschoon/gofit/chart"
	"github.com/kevinschoon/gofit/database"
	"github.com/kevinschoon/gofit/models"

	"net/http"
	"text/template"
)

const (
	staticDir       string = "www"
	baseTmpl        string = staticDir + "/base.html"
	chartTmpl       string = staticDir + "/chart.html"
	dataTmpl        string = staticDir + "/data.html"
	panelTmpl       string = staticDir + "/panel.html"
	collectionsTmpl string = staticDir + "/collections.html"
)

type Response struct {
	Title       string
	Query       *Query
	Collection  *models.Collection
	Collections []string
	URLBuilder  *URLBuilder
	Templates   []string
}

type Handler struct {
	db        *database.DB
	handle    func(http.ResponseWriter, *http.Request, *Response) error
	Templates []string
}

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := &Response{
		Query: QueryFromURL(r.URL),
		URLBuilder: &URLBuilder{
			URL: r.URL,
		},
		Templates: handler.Templates,
	}
	collections, err := handler.db.Collections()
	if err != nil {
		handler.Err(err, w, r)
		return
	}
	response.Collections = collections
	if name, ok := mux.Vars(r)["collection"]; ok {
		response.URLBuilder.Collection = name
		collection, err := handler.db.Read(name, response.Query.Start, response.Query.End)
		if err != nil {
			handler.Err(err, w, r)
			return
		}
		collection.RollUp(response.Query.Precision)
		response.Collection = collection
	}
	handler.Err(handler.handle(w, r, response), w, r)
}

func (handler Handler) Err(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		switch err.(type) {
		case template.ExecError:
		default:
			switch err.Error() {
			case database.ErrCollectionNotFound.Error():
				http.NotFound(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func Chart(w http.ResponseWriter, r *http.Request, response *Response) error {
	canvas, err := chart.New(response.Collection, response.Query.X, response.Query.Y)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("Vary", "Accept-Encoding")
	_, err = canvas.WriteTo(w)
	return err
}

func RenderHome(w http.ResponseWriter, r *http.Request, response *Response) error {
	tmpl, err := template.ParseFiles(response.Templates...)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, response)
}
