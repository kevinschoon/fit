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

type TemplateData struct {
	Activities tcx.Acts
	Activity   tcx.Activity
	ChartURL   string
	Lap        tcx.Lap
}

func StartEnd(r *http.Request) (start, end time.Time, err error) {
	query := r.URL.Query()
	s, e := query.Get("start"), query.Get("end")
	if s != "" && e != "" {
		start, err = time.Parse(qTime, s)
		if err != nil {
			return start, end, err
		}
		end, err = time.Parse(qTime, e)
		if err != nil {
			return start, end, err
		}
	} else {
		end = time.Now()
		start = time.Date(end.Year()-1, // Show last year of data by default
			end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
	}
	return start, end, err
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
	start, end, err := StartEnd(r)
	if err != nil {
		return err
	}
	activities, err := Activities(db, Between(start, end, "start_time"))
	if err != nil {
		return err
	}
	return tmpl.Execute(w, &TemplateData{
		Activities: activities,
		ChartURL:   fmt.Sprintf("/chart?%s", r.URL.RawQuery),
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

func HandleChart(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	start, end, err := StartEnd(r)
	if err != nil {
		return err
	}
	activities, err := Activities(db, Between(start, end, "start_time"))
	if err != nil {
		return err
	}
	canvas, err := DistanceOverTime(activities)
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
	router.Handle("/chart", BasicHandler(HandleChart))
	router.Handle("/chart/{activity}", BasicHandler(HandleDetailsChart))
	router.Handle("/chart/{activity}/lap/{lap}", BasicHandler(HandleDetailsChart))
	log.Printf("Fit server listening @ %s", listenPattern)
	log.Fatal(http.ListenAndServe(listenPattern, handlers.CombinedLoggingHandler(os.Stdout, router)))
}
