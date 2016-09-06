package server

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kevinschoon/gofit/chart"
	"github.com/kevinschoon/gofit/database"
	"github.com/kevinschoon/gofit/models"
	"github.com/kevinschoon/gofit/models/tcx"

	"net/http"
	"text/template"
)

type SeriesHandler struct {
	dbPath string
	handle func(models.Series, http.ResponseWriter, *http.Request) error
}

func (handler SeriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	loader := tcx.TCXLoader{} // TODO Switch out for dynamic types
	db, err := database.New(handler.dbPath, loader)
	if err != nil {
		handler.HandleError(nil, err, w)
		return
	}
	defer db.Close()
	query := models.QueryFromURL(r.URL)
	series, err := database.Read(db, query, loader)
	if err != nil {
		handler.HandleError(nil, err, w)
		return
	}
	handler.HandleError(db, handler.handle(series, w, r), w)
}

func (handler SeriesHandler) HandleError(db *gorm.DB, err error, w http.ResponseWriter) {
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		switch err.(type) {
		case template.ExecError:
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

func HandleActivities(series models.Series, w http.ResponseWriter, r *http.Request) error {
	tmpl, data, err := LoadTemplate(r)
	if err != nil {
		return err
	}
	data.Columns = series.Columns()
	data.Rows = series.Rows()
	return tmpl.Execute(w, data)
}

func HandleChart(series models.Series, w http.ResponseWriter, r *http.Request) error {
	canvas, err := chart.New(series.Pts(""))
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("Vary", "Accept-Encoding")
	_, err = canvas.WriteTo(w)
	return err
}
