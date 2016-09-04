package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/kevinschoon/tcx"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	text "text/template"
	"time"
)

const (
	StaticDir string = "www"
	baseTmpl  string = StaticDir + "/base.html"
	chartTmpl string = StaticDir + "/chart.html"
	dataTmpl  string = StaticDir + "/data.html"
)

func ThisWeek(now time.Time) time.Time {
	return now.AddDate(0, 0, -7)
}

func ThisMonth(now time.Time) time.Time {
	return now.AddDate(0, -1, 0)
}

func ThisYear(now time.Time) time.Time {
	return now.AddDate(-1, 0, 0)
}

func CopyURL(u *url.URL) url.URL {
	return url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
	}
}

type Query struct {
	Sport string
	Aggr  string
	URL   *url.URL
	Start time.Time
	End   time.Time
}

func (q Query) ByRange(key string) string {
	u := CopyURL(q.URL)
	values := q.URL.Query()
	now := time.Now()
	values.Set("end", now.Format(qTime))
	switch key {
	case "week":
		values.Set("start", ThisWeek(now).Format(qTime))
	case "month":
		values.Set("start", ThisMonth(now).Format(qTime))
	case "year":
		values.Set("start", ThisYear(now).Format(qTime))
	}
	u.RawQuery = values.Encode()
	return u.String()
}

func (q Query) ByTime(t time.Time) string {
	u := CopyURL(q.URL)
	values := q.URL.Query()
	values.Set("end", t.Format(qTime))
	values.Set("start", t.Format(qTime))
	values.Set("aggr", "none")
	u.RawQuery = values.Encode()
	return u.String()
}

func (q Query) BySport(key string) string {
	u := CopyURL(q.URL)
	values := q.URL.Query()
	values.Set("sport", key)
	if key == "all" {
		values.Del("sport")
	}
	u.RawQuery = values.Encode()
	return u.String()
}

func (q Query) ByAggr(key string) string {
	u := CopyURL(q.URL)
	values := q.URL.Query()
	values.Set("aggr", key)
	u.RawQuery = values.Encode()
	return u.String()
}

// NewQuery returns a query object from an HTTP request
func NewQuery(request *http.Request) (query Query) {
	query.URL = request.URL
	query.Start, _ = time.Parse(qTime, query.URL.Query().Get("start"))
	query.End, _ = time.Parse(qTime, query.URL.Query().Get("end"))
	if query.Start.IsZero() || query.End.IsZero() {
		query.End = time.Now()
		query.Start = ThisMonth(query.End)
	}
	query.Sport = query.URL.Query().Get("sport")
	if query.Sport == "" {
		query.Sport = "%"
	}
	query.Aggr = query.URL.Query().Get("aggr")
	if query.Aggr == "" {
		query.Aggr = "day"
	}
	return query
}

type TemplateData struct {
	Activities tcx.Acts
	Activity   tcx.Activity
	ChartURL   string
	Lap        tcx.Lap
	Sports     []string
	Query      Query
}

func LoadTemplates() (*template.Template, error) {
	return template.ParseFiles(baseTmpl, chartTmpl, dataTmpl)
}

type BasicHandler func(*gorm.DB, http.ResponseWriter, *http.Request) error

func (fn BasicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db, err := GetDB()
	if err != nil {
		HandleError(nil, err, w)
		return
	}
	defer db.Close()
	HandleError(db, fn(db, w, r), w)
}

func HandleError(db *gorm.DB, err error, w http.ResponseWriter) {
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		switch err.(type) {
		case text.ExecError:
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		if db != nil {
			if db.Error != nil {
				fmt.Println("DB Error:", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func HandleActivities(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	tmpl, err := LoadTemplates()
	if err != nil {
		return err
	}
	query := NewQuery(r)
	sports, err := Sports(db)
	if err != nil {
		return err
	}
	activities, err := Activities(db, FromQuery(query))
	if err != nil {
		return err
	}
	activities = RollUpActivities(activities, query.Aggr)
	return tmpl.Execute(w, &TemplateData{
		Activities: activities,
		ChartURL:   fmt.Sprintf("/chart?%s", r.URL.RawQuery),
		Sports:     sports,
		Query:      query,
	})
}

func HandleChart(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	query := NewQuery(r)
	activities, err := Activities(db, FromQuery(query))
	if err != nil {
		return err
	}
	canvas, err := DistanceOverTime(ActivityByDist(RollUpActivities(activities, query.Aggr)))
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("Vary", "Accept-Encoding")
	_, err = canvas.WriteTo(w)
	return err
}

func RunServer(listenPattern string) {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", BasicHandler(HandleActivities))
	router.HandleFunc("/static/dashboard.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, StaticDir+"/dashboard.css")
	})
	router.Handle("/chart", BasicHandler(HandleChart))
	log.Printf("Fit server listening @ %s", listenPattern)
	log.Fatal(http.ListenAndServe(listenPattern, handlers.CombinedLoggingHandler(os.Stdout, router)))
}
