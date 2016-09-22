package server

import (
	"fmt"
	"time"

	"github.com/gorilla/mux"
	"github.com/kevinschoon/fit/chart"
	"github.com/kevinschoon/fit/database"
	"github.com/kevinschoon/fit/models"

	"net/http"
	"net/url"
	"text/template"
)

type Response struct {
	Title    string
	Series   []*models.Series
	Choices  []*models.Series
	Explore  bool     // Display Data Explorer
	Browse   bool     // Display Series Listing
	Keys     []string // Series Keys to Display
	ChartURL string   // URL for rendering the chart
	Query    url.Values
	DemoMode bool
	Version  string
}

// Match checks if the a query string is set to a particular value
func (r Response) Match(name, value string) bool { return r.Query.Get(name) == value }

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
			case database.ErrSeriesNotFound:
				http.NotFound(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

type Handler struct {
	db        *database.DB
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
	return err
}

func (handler Handler) Home(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles(handler.templates...)
	if err != nil {
		return err
	}
	start, end := StartEnd(r.URL)
	response := handler.response()
	choices, err := handler.db.Series()
	if err != nil {
		return err
	}
	response.Choices = choices
	response.Query = r.URL.Query()
	if name, ok := mux.Vars(r)["series"]; ok {
		response.Title = name
		response.Explore = true
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
	}
	response.Title = "Browse"
	response.Browse = true
	return tmpl.Execute(w, response)
}
