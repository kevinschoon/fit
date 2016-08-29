package main

import (
	"github.com/jinzhu/gorm"
	"github.com/kevinschoon/tcx"
	"net/http"
	"time"
)

const (
	Paginate int    = 100
	dbTime   string = "2006-01-02"  // Date format for SQL query
	qTime    string = "2006-Jan-02" // Date format URL query
)

// Writer writes activitiy data into the database
type Writer interface {
	Write(*gorm.DB, []tcx.Activity) error
}

// Reader reads activity data from the database
type Reader interface {
	Read(*gorm.DB) ([]tcx.Activity, error)
}

// Simple query reads activities within a given range
type WebQuery struct {
	Start time.Time // return entries after Start
	End   time.Time // return entries before End
	Limit int       // Maximum batch size of entries
	Query string    // Raw SQL query string
}

func (query WebQuery) query(db *gorm.DB) *gorm.DB {
	return db.Where(query.Query, query.Start.Format(dbTime), query.End.Format(dbTime))
}

func (query WebQuery) Read(db *gorm.DB) (activities []tcx.Activity, err error) {
	var (
		count int
		last  int
	)
	query.query(db.Model(&tcx.Activity{})).Count(&count)
	for len(activities) < count { // TODO: Cleanup
		if db.Error != nil {
			return nil, db.Error
		}
		var results []tcx.Activity
		query.query(db).Limit(query.Limit).Offset(last).Preload("Laps.Trk").Find(&results)
		for _, result := range results {
			activities = append(activities, result)
		}
		last += query.Limit
	}
	//db.Preload("Laps.Trk.Pt").Find(&activities)
	return activities, nil
}

// NewSimpleQuery returns a SimpleQuery Reader interface
func NewWebQuery(r *http.Request) (Reader, error) {
	q := r.URL.Query()
	start, end := q.Get("start"), q.Get("end")
	query := WebQuery{
		Limit: Paginate,
		Query: "start_time >= ? AND start_time <= ?",
	}
	if start != "" && end != "" {
		s, err := time.Parse(qTime, start)
		if err != nil {
			return nil, err
		}
		query.Start = s
		e, err := time.Parse(qTime, end)
		if err != nil {
			return nil, err
		}
		query.End = e
	} else { // Default to 7 days prior
		query.End = time.Now()
		query.Start = time.Date(query.End.Year(),
			query.End.Month(), query.End.Day()-7, 0, 0, 0, 0, time.UTC)
	}
	return query, nil
}
