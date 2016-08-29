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
)

const StaticDir string = "www"

type TemplateData struct {
	Activities []tcx.Activity
	Activity   tcx.Activity
	QueryStr   string
	Lap        tcx.Lap
}

type ChartNotFound struct {
	Name string
}

func (cnf ChartNotFound) Error() string {
	return fmt.Sprintf("Chart %s not found", cnf.Name)
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
		case ChartNotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
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
	query, err := NewWebQuery(r)
	if err != nil {
		return err
	}
	activities, err := query.Read(db)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, &TemplateData{
		QueryStr:   r.URL.RawQuery,
		Activities: activities,
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
	}
	if len(activity.Laps) == 1 {
		data.Lap = activity.Laps[0]
	}
	if mux.Vars(r)["lap"] != "" {
		lapId, err := strconv.ParseInt(mux.Vars(r)["lap"], 10, 64)
		if err != nil {
			return err
		}
		fmt.Println(len(activity.Laps), int(lapId))
		if !(len(activity.Laps) >= int(lapId)) {
			return fmt.Errorf("Lap Not Found")
		}
		data.Lap = activity.Laps[int(lapId)]
	}
	return tmpl.Execute(w, data)
}

/*
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
*/

func RunServer(listenPattern string) {
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", BasicHandler(HandleActivities))
	router.Handle("/activity/{activity}", BasicHandler(HandleActivity))
	router.Handle("/activity/{activity}/lap/{lap}", BasicHandler(HandleActivity))
	router.HandleFunc("/static/dashboard.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, StaticDir+"/dashboard.css")
	})
	//router.Handle("/chart/{name}", BasicHandler(HandleGraph))
	log.Printf("Fit server listening @ %s", listenPattern)
	log.Fatal(http.ListenAndServe(listenPattern, handlers.CombinedLoggingHandler(os.Stdout, router)))
}
