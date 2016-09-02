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
	"os"
	"strconv"
	text "text/template"
	"time"
)

const StaticDir string = "www"

type Query struct {
	Sport string
	Last  string
	Start time.Time
	End   time.Time
}

func (q Query) Where(table string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch table {
		case "activities":
			return db.Where("start_time >= ? AND start_time <= ? AND sport LIKE ?", q.Start.Format(dbTime), q.End.Format(dbTime), q.Sport)
		default:
			return db.Where("time >= ? AND time <= ?", q.Start.Format(dbTime), q.End.Format(dbTime))
		}
	}
}

// NewQuery returns a query object from an HTTP request
func NewQuery(r *http.Request) (query Query) {
	values := r.URL.Query()
	query.Start, _ = time.Parse(qTime, values.Get("start"))
	query.End, _ = time.Parse(qTime, values.Get("end"))
	if query.Start.IsZero() || query.End.IsZero() {
		query.End = time.Now()
		query.Start = time.Date(query.End.Year(), query.End.Month(), query.End.Day()-7, 0, 0, 0, 0, time.UTC)
		query.Last = values.Get("last")
		if query.Last != "" {
			switch query.Last {
			case "week":
			case "month":
				query.Start = time.Date(query.End.Year(), query.End.Month()-1, query.End.Day(), 0, 0, 0, 0, time.UTC)
			case "year":
				query.Start = time.Date(query.End.Year()-1, query.End.Month(), query.End.Day(), 0, 0, 0, 0, time.UTC)
			}
		}
	}
	query.Sport = values.Get("sport")
	if query.Sport == "" {
		query.Sport = "%"
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

func LoadTemplates(section string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/%s", StaticDir, "base.html"))
	if err != nil {
		return nil, err
	}
	tmpl, err = tmpl.ParseFiles(fmt.Sprintf("%s/%s", StaticDir, section))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
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
	tmpl, err := LoadTemplates("activities.html")
	if err != nil {
		return err
	}
	query := NewQuery(r)
	sports, err := Sports(db)
	if err != nil {
		return err
	}
	activities, err := Activities(db, query.Where("activities"))
	if err != nil {
		return err
	}
	activities = RollUpActivities(activities, "day")
	return tmpl.Execute(w, &TemplateData{
		Activities: activities,
		ChartURL:   fmt.Sprintf("/chart?%s", r.URL.RawQuery),
		Sports:     sports,
	})
}

func HandleActivity(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	tmpl, err := LoadTemplates("activity.html")
	if err != nil {
		return err
	}
	activityId, err := strconv.ParseUint(mux.Vars(r)["activity"], 10, 64)
	if err != nil {
		return err
	}
	activity := Activity(uint(activityId), db)
	data := &TemplateData{
		Activity: activity,
		ChartURL: fmt.Sprintf("/chart/%d", activityId),
	}
	if len(activity.Laps) == 1 {
		data.Lap = activity.Laps[0]
	}
	if mux.Vars(r)["lap"] != "" {
		lapId, err := strconv.ParseInt(mux.Vars(r)["lap"], 10, 64)
		if err != nil {
			return err
		}
		if !(len(activity.Laps) >= int(lapId)) {
			return fmt.Errorf("Lap Not Found")
		}
		data.ChartURL += fmt.Sprintf("/lap/%d", lapId)
		data.Lap = activity.Laps[int(lapId)]
	}
	return tmpl.Execute(w, data)
}

func HandleTrackPoints(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	query := NewQuery(r)
	pts := Trackpoints(db, query.Where("trackpoints"))
	fmt.Println(pts)
	return nil
}

func HandleChart(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	query := NewQuery(r)
	activities, err := Activities(db, query.Where("activities"))
	if err != nil {
		return err
	}
	canvas, err := DistanceOverTime(ActivityByDist(RollUpActivities(activities, "day")))
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "image/svg+xml")
	w.Header().Add("Vary", "Accept-Encoding")
	_, err = canvas.WriteTo(w)
	return err
}

func HandleDetailsChart(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	var lap tcx.Lap
	activityId, err := strconv.ParseUint(mux.Vars(r)["activity"], 10, 64)
	if err != nil {
		return err
	}
	activity := Activity(uint(activityId), db)
	if len(activity.Laps) == 1 {
		lap = activity.Laps[0]
	}
	if mux.Vars(r)["lap"] != "" {
		lapId, err := strconv.ParseInt(mux.Vars(r)["lap"], 10, 64)
		if err != nil {
			return err
		}
		if !(len(activity.Laps) >= int(lapId)) {
			return fmt.Errorf("Lap Not Found")
		}
		lap = activity.Laps[int(lapId)]
	}
	canvas, err := ChartXYs(tcx.Trackpoints(lap.Trk.Pt))
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "image/svg+xml")
	w.Header().Add("Vary", "Accept-Encoding")
	_, err = canvas.WriteTo(w)
	return err
}

func RunServer(listenPattern string) {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", BasicHandler(HandleActivities))
	router.Handle("/activity/{activity}", BasicHandler(HandleActivity))
	router.Handle("/activity/{activity}/lap/{lap}", BasicHandler(HandleActivity))
	router.HandleFunc("/static/dashboard.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, StaticDir+"/dashboard.css")
	})
	router.Handle("/data", BasicHandler(HandleTrackPoints))
	router.Handle("/chart", BasicHandler(HandleChart))
	router.Handle("/chart/{activity}", BasicHandler(HandleDetailsChart))
	router.Handle("/chart/{activity}/lap/{lap}", BasicHandler(HandleDetailsChart))
	log.Printf("Fit server listening @ %s", listenPattern)
	log.Fatal(http.ListenAndServe(listenPattern, handlers.CombinedLoggingHandler(os.Stdout, router)))
}
