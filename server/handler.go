package server

import (
	"fmt"
	"github.com/gorilla/mux"
	//"github.com/kevinschoon/gofit/chart"
	"github.com/kevinschoon/gofit/database"
	"github.com/kevinschoon/gofit/models"
	"time"

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
	Explore bool     // Display Data Explorer
	Browse  bool     // Display Series Listing
	Keys    []string // Series Keys to Display
}

// Rows iterates each series for templating
func (r Response) Rows() [][]string {
	rows := make([][]string, len(r.Series))
	for s, series := range r.Series {
		rows[s] = make([]string, len(r.Keys))
		for k, key := range r.Keys {
			// Only show a single level of value since these
			// series should already be aggregated
			rows[s][k] = series.Value(0, key).String()
			if key == "time" {
				rows[s][k] = series.Value(0, key).Time().Format(time.RFC3339)
			}
		}
	}
	return rows
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
	start, end := StartEnd(r.URL)
	series, err := handler.db.ReadSeries(mux.Vars(r)["series"], start, end)
	if err != nil {
		return err
	}
	fmt.Println(len(series))
	series = models.Resize(series, Aggr(r.URL))
	fmt.Println(series[0].Dump())
	series = models.Apply(series, Fn(r.URL))
	fmt.Println(series[0].Dump())
	return nil
}

func (handler Handler) Home(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles(handler.templates...)
	if err != nil {
		return err
	}
	start, end := StartEnd(r.URL)
	response := &Response{}
	if name, ok := mux.Vars(r)["series"]; ok {
		response.Title = "Explorer"
		response.Explore = true
		series, err := handler.db.ReadSeries(name, start, end)
		if err != nil {
			return err
		}
		response.Keys = StrArray("keys", r.URL)
		if len(response.Keys) < 1 {
			response.Keys = models.Keys(series[0])
		}
		series = models.Resize(series, Aggr(r.URL))
		series = models.Apply(series, Fn(r.URL))
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
