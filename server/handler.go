package server

import (
	"fmt"
	//"time"

	//"github.com/gorilla/mux"
	//"github.com/kevinschoon/fit/chart"
	//mtx "github.com/gonum/matrix/mat64"
	"github.com/kevinschoon/fit/store"

	"net/http"
	"net/url"
	"text/template"
)

type Response struct {
	Title    string
	Explore  bool     // Display Data Explorer
	Browse   bool     // Display Series Listing
	Keys     []string // Series Keys to Display
	ChartURL string   // URL for rendering the chart
	Datasets []*store.Dataset
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
	/*
		start, end := StartEnd(r.URL)
		series, err := handler.db.ReadSeries(mux.Vars(r)["series"], start, end)
		if err != nil {
			return err
		}
		series = models.Resize(series, Aggr(r.URL))
		series = models.Apply(series, Fn(r.URL))
		canvas, err := chart.New(ChartCfg(series[0], r.URL), series)
		if err != nil {
			return err
		}
		_, err = canvas.WriteTo(w)
	*/
	return nil
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
	/*
		response.ChartURL = Chart(r.URL)
		series, err := handler.db.ReadSeries(name, start, end)
		if err != nil {
			return err
		}
		response.Keys = StrArray("keys", r.URL)
		if len(response.Keys) < 1 {
			response.Keys = series[0].Keys.Names()
		}
		series = models.Resize(series, Aggr(r.URL))
		series = models.Apply(series, Fn(r.URL))
		response.Series = series
		return tmpl.Execute(w, response)
		response.Title = "Browse"
		response.Browse = true
	*/
	return tmpl.Execute(w, response)
}
