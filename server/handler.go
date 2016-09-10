package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kevinschoon/gofit/chart"
	"github.com/kevinschoon/gofit/database"
	"github.com/kevinschoon/gofit/models"

	"net/http"
	"text/template"
)

type CollectionHandler struct {
	db     *database.DB
	handle func(*models.Collection, models.Query, http.ResponseWriter, *http.Request) error
}

func (handler CollectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := models.QueryFromURL(r.URL)
	collection, err := handler.db.Read(mux.Vars(r)["collection"], query.Start, query.End)
	collection.RollUp(query.Precision)
	if err != nil {
		handler.HandleError(err, w, r)
		return
	}
	handler.HandleError(handler.handle(collection, query, w, r), w, r)
}

func (handler CollectionHandler) HandleError(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		switch err.(type) {
		case template.ExecError:
		default:
			switch err.Error() {
			case database.ErrCollectionNotFound.Error():
				http.NotFound(w, r)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func Collection(collection *models.Collection, query models.Query, w http.ResponseWriter, r *http.Request) error {
	tmpl, data, err := LoadTemplate(r, collection)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

func Chart(collection *models.Collection, query models.Query, w http.ResponseWriter, r *http.Request) error {
	canvas, err := chart.New(collection)
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "image/png")
	w.Header().Add("Vary", "Accept-Encoding")
	_, err = canvas.WriteTo(w)
	return nil
}
