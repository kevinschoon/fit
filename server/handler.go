package server

import (
	"encoding/json"
	"fmt"
	"github.com/gonum/plot/vg"
	"github.com/kevinschoon/fit/chart"
	"github.com/kevinschoon/fit/types"
	"image/color"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
)

type Response struct {
	Title    string
	Explore  bool     // Display Data Explorer
	Browse   bool     // Display Datasets Listing
	Keys     []string // Series Keys to Display
	ChartURL string   // URL for rendering the chart
	Datasets []*types.Dataset
	Dataset  *types.Dataset
	Query    url.Values
	DemoMode bool
	Version  string
}

type Handler struct {
	db        types.Client
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
	ds, err := handler.db.Query(types.NewQueryFromQS(r.URL))
	if err != nil {
		return err
	}
	cfg := chart.Config{
		Title:          ds.Name,
		PrimaryColor:   color.White,
		SecondaryColor: color.Black,
		Width:          18 * vg.Inch,
		Height:         5 * vg.Inch,
		Columns:        types.NewQueryFromQS(r.URL).Columns(),
	}
	if w, err := strconv.ParseInt(r.URL.Query().Get("width"), 0, 64); err == nil {
		if w < 20 { // Prevent potentially horrible DOS
			cfg.Width = vg.Length(w) * vg.Inch
		}
	}
	if h, err := strconv.ParseInt(r.URL.Query().Get("height"), 0, 64); err == nil {
		if h < 20 {
			cfg.Height = vg.Length(h) * vg.Inch
		}
	}
	canvas, err := chart.New(cfg, ds.Mtx)
	if err != nil {
		return err
	}
	_, err = canvas.WriteTo(w)
	return err
}

func (handler Handler) DatasetAPI(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		query := types.NewQueryFromQS(r.URL)
		fmt.Println(query.Col, query.Max, query.Function)
		if query.Len() > 0 { // If URL contains a query return the query result
			ds, err := handler.db.Query(query)
			if err != nil {
				return err
			}
			ds.WithValues = true // Encode values in the dataset response
			if err := json.NewEncoder(w).Encode(ds); err != nil {
				return err
			}
		} else { // Otherwise return an array of all datasets
			datasets, err := handler.db.Datasets()
			if err != nil {
				return err
			}
			if err = json.NewEncoder(w).Encode(datasets); err != nil {
				return err
			}
		}
	case "POST":
		ds := &types.Dataset{WithValues: true}
		if err := json.NewDecoder(r.Body).Decode(ds); err != nil {
			return err
		}
		if err := handler.db.Write(ds); err != nil {
			return err
		}
	case "DELETE":
		if name := r.URL.Query().Get("name"); name != "" {
			if err := handler.db.Delete(name); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("specify name")
		}
	}
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
	chartURL := &url.URL{
		Scheme:   r.URL.Scheme,
		Host:     r.URL.Host,
		Path:     "/chart",
		RawQuery: r.URL.Query().Encode(),
	}
	response.ChartURL = chartURL.String()
	ds, err := handler.db.Query(types.NewQueryFromQS(r.URL))
	if err != nil {
		return err
	}
	response.Dataset = ds
	return tmpl.Execute(w, response)
}
