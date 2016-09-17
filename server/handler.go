package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"time"
	//"github.com/kevinschoon/gofit/chart"
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
	Title   string
	Series  []*models.Series
	Explore bool // Display Data Explorer
	Browse  bool // Display Series Listing
}

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		switch err.(type) {
		case template.ExecError:
		default:
			switch err {
			case database.ErrSeriesNotFound:
				http.NotFound(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

type ErrorHandler func(http.ResponseWriter, *http.Request) error

func (fn ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleError(fn(w, r), w, r)
}

type Handler struct {
	db        *database.DB
	templates []string
}

func (handler Handler) Chart(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (handler Handler) Home(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles(handler.templates...)
	if err != nil {
		return err
	}
	response := &Response{}
	if name, ok := mux.Vars(r)["series"]; ok {
		response.Title = "Explorer"
		response.Explore = true
		series, err := handler.db.ReadSeries(name, time.Time{}, time.Now())
		if err != nil {
			return err
		}
		response.Series = series
		return tmpl.Execute(w, response)
	}
	response.Title = "Browse"
	response.Browse = true
	series, err := handler.db.Series()
	if err != nil {
		return err
	}
	response.Series = series
	return tmpl.Execute(w, response)
}
