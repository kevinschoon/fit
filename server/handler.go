package server

import (
	"fmt"
	"github.com/kevinschoon/fit/chart"
	"github.com/kevinschoon/fit/store"
	"net/http"
	"net/url"
	"text/template"
)

type Response struct {
	Title    string
	Explore  bool     // Display Data Explorer
	Browse   bool     // Display Datasets Listing
	Keys     []string // Series Keys to Display
	ChartURL string   // URL for rendering the chart
	Datasets []*store.Dataset
	Dataset  *store.Dataset
	Query    url.Values
	DemoMode bool
	Version  string
}

type ErrorHandler func(http.ResponseWriter, *http.Request) error

func (fn ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleError(fn(w, r), w, r)
}

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		switch err.(type) {
		case template.ExecError:
		default:
			switch err {
			case store.ErrNotFound:
				http.NotFound(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

type Handler struct {
	db        *store.DB
	version   string
	templates []string
	defaults  Response
}

func (handler Handler) response() *Response {
	return &Response{
		DemoMode: handler.defaults.DemoMode,
		Version:  handler.defaults.Version,
	}
}

func (handler Handler) Chart(w http.ResponseWriter, r *http.Request) error {
	ds, err := handler.db.Query(XYQueries(r.URL))
	if err != nil {
		return err
	}
	cfg := ChartCfg(ds, r.URL)
	canvas, err := chart.New(cfg, ds.Mtx)
	if err != nil {
		return err
	}
	_, err = canvas.WriteTo(w)
	return err
}

func (handler Handler) Explore(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles(handler.templates...)
	if err != nil {
		return err
	}
	response := handler.response()
	datasets, err := handler.db.Datasets()
	if err != nil {
		return err
	}
	response.Datasets = datasets
	response.Query = r.URL.Query()
	response.Explore = true
	response.ChartURL = Chart(r.URL)
	ds, err := handler.db.Query(XYQueries(r.URL))
	if err != nil {
		return err
	}
	response.Dataset = ds
	return tmpl.Execute(w, response)
}

func (handler Handler) Home(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles(handler.templates...)
	if err != nil {
		return err
	}
	response := handler.response()
	datasets, err := handler.db.Datasets()
	if err != nil {
		return err
	}
	response.Datasets = datasets
	response.Query = r.URL.Query()
	response.Title = "Browse"
	response.Browse = true
	return tmpl.Execute(w, response)
}
