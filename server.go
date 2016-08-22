package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const timeStr string = "2006-Jan-02"

type TemplateData struct {
	Section         string
	OverviewURL     string
	RegressionURL   string
	DistributionURL string
	Bucket          string
	Total           *Total
	Totals          []*Total
}

type ChartNotFound struct {
	Name string
}

func (cnf ChartNotFound) Error() string {
	return fmt.Sprintf("Chart %s not found", cnf.Name)
}

type BasicHandler func(http.ResponseWriter, *http.Request) error

func (fn BasicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleError(fn(w, r), w)
}

func HandleError(err error, w http.ResponseWriter) {
	if err != nil {
		switch err.(type) {
		case ChartNotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetFilter(query url.Values) (Filter, error) {
	var (
		start  string
		end    string
		filter Filter
	)
	start = query.Get("start")
	end = query.Get("end")
	if start != "" && end != "" {
		s, err := time.Parse(timeStr, start)
		if err != nil {
			return filter, err
		}
		e, err := time.Parse(timeStr, end)
		if err != nil {
			return filter, err
		}
		filter = Filter(DateFilter(s, e))
	} else {
		filter = Filter(NullFilter)
	}
	return filter, nil
}

func GetBucket(query url.Values) string {
	if query.Get("bucket") == "" {
		return "month"
	}
	return query.Get("bucket")
}

func HandleHome(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	filter, err := GetFilter(query)
	if err != nil {
		return err
	}
	bucket := GetBucket(query)
	tmpl := template.Must(template.ParseFiles(StaticDir + "/index.html"))
	return tmpl.Execute(w, &TemplateData{
		Section:         "Overview",
		Bucket:          bucket,
		Total:           database.total,
		Totals:          database.Totals(filter, bucket),
		OverviewURL:     fmt.Sprintf("chart/overview?%s", query.Encode()),
		RegressionURL:   fmt.Sprintf("chart/regression?%s", query.Encode()),
		DistributionURL: fmt.Sprintf("chart/distribution?%s", query.Encode()),
	})
}

func HandleGraph(w http.ResponseWriter, r *http.Request) error {
	filter, err := GetFilter(r.URL.Query())
	if err != nil {
		return err
	}
	name := mux.Vars(r)["name"]
	var chart Chart
	switch name {
	case "overview":
		chart = OverviewChart{Title: "Distance By Totals"}
	case "regression":
		chart = RegressionChart{}
	case "distribution":
		chart = DistributionChart{}
	default:
		return ChartNotFound{Name: name}
	}
	w.Header().Add("Content-Type", "image/svg+xml")
	w.Header().Add("Vary", "Accept-Encoding")
	canvas, err := chart.Canvas(database.Totals(filter, GetBucket(r.URL.Query())))
	if err != nil {
		return err
	}
	_, err = canvas.WriteTo(w)
	return err
}

func RunServer(listenPattern string) {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", BasicHandler(HandleHome))
	router.Handle("/{name}", BasicHandler(HandleHome))
	router.HandleFunc("/static/dashboard.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, StaticDir+"/dashboard.css")
	})
	router.Handle("/chart/{name}", BasicHandler(HandleGraph))
	log.Printf("Fit server listening @ %s", listenPattern)
	log.Fatal(http.ListenAndServe(listenPattern, handlers.CombinedLoggingHandler(os.Stdout, router)))
}
