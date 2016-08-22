package main

import (
	"fmt"
	"github.com/philhofer/tcx"
	"os"
	"path/filepath"
	"time"
)

// database is a singleton
var database *Database

// Filter is used to filter loaded tcx data in various ways
type Filter func(lap tcx.Lap) bool

// NullFilter always returns true
func NullFilter(lap tcx.Lap) bool { return true }

// DateFilter returns a Filter which is true if the given tcx.Lap
// falls between the provided start and end times
func DateFilter(start, end time.Time) func(lap tcx.Lap) bool {
	return func(lap tcx.Lap) bool {
		if lap.Start.After(start) && lap.Start.Add(time.Second*time.Duration(lap.TotalTime)).Before(end) {
			return true
		}
		return false
	}
}

// Total represents aggregated tcx lap and track data
type Total struct {
	Start     time.Time
	Bucket    string
	TotalTime float64
	Dist      float64
}

func (t Total) Km() float64 {
	return t.Dist * 0.001
}

func (t Total) Name() (name string) {
	switch t.Bucket {
	case "year":
		name = fmt.Sprintf("%d", t.Start.Year())
	case "month":
		name = fmt.Sprintf("%d", int(t.Start.Month()))
	case "day":
		name = fmt.Sprintf("%d", t.Start.Day())
	}
	return name
}

// Totals is an array of Total
type Totals []*Total

// Database is an singleton containing tcx lap data
type Database struct {
	laps  tcx.Laps
	total *Total
}

// Laps returns all tcx.Lap for which the given filter returns true
func (db *Database) Laps(fn Filter) []*tcx.Lap {
	laps := make([]*tcx.Lap, 0)
	for i := 0; i < len(db.laps); i++ {
		if fn(db.laps[i]) {
			laps = append(laps, &db.laps[i])
		}
	}
	return laps
}

// Total aggregats all tcx laps together
func (db *Database) Total(fn Filter) *Total {
	total := &Total{}
	if len(db.laps) > 0 {
		total.Start = db.laps[0].Start
	}
	for _, lap := range db.Laps(fn) {
		total.Dist += lap.Dist
		total.TotalTime += lap.TotalTime
	}
	return total
}

// Totals aggregates tcx laps by "year", "month", or "day"
func (db *Database) Totals(fn Filter, bucket string) []*Total {
	t := make(map[int]*Total)
	totals := make(Totals, 0)
	for _, lap := range db.Laps(fn) {
		switch bucket {
		case "year":
			key := TimeKey(lap.Start, bucket)
			if _, ok := t[key]; !ok {
				t[key] = &Total{
					Start: time.Date(lap.Start.Year(),
						time.January, 01, 0, 0, 0, 0, time.UTC),
					Bucket: bucket,
				}
				totals = append(totals, t[key])
			}
			t[key].Dist += lap.Dist
			t[key].TotalTime += lap.TotalTime
		case "month":
			key := TimeKey(lap.Start, bucket)
			if _, ok := t[key]; !ok {
				t[key] = &Total{
					Start: time.Date(lap.Start.Year(),
						lap.Start.Month(), 01, 0, 0, 0, 0, time.UTC),
					Bucket: bucket,
				}
				totals = append(totals, t[key])
			}
			t[key].Dist += lap.Dist
			t[key].TotalTime += lap.TotalTime
		case "day":
			key := TimeKey(lap.Start, bucket)
			if _, ok := t[key]; !ok {
				t[key] = &Total{
					Start: time.Date(lap.Start.Year(),
						lap.Start.Month(), lap.Start.Day(), 0, 0, 0, 0, time.UTC),
					Bucket: bucket,
				}
				totals = append(totals, t[key])
			}
			t[key].Dist += lap.Dist
			t[key].TotalTime += lap.TotalTime
		}
	}
	return totals
}

// Generate a unique key for a given date
func TimeKey(t time.Time, precision string) (key int) {
	switch precision {
	case "year":
		key = t.Year()
	case "month":
		key = (t.Year() * 12) + int(t.Month())
	case "day":
		key = (t.Year() * 12) + (int(t.Month()) * 31) + t.Day()
	}
	return key
}

// InitDB initializes the global "database" struct when provided
// a directory of files containing TCX data. It has been tested
// with TCX data provided by Google Checkout.
func InitDB(directory string) (err error) {
	database = &Database{
		laps: make(tcx.Laps, 0),
	}
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			tdb, err := tcx.ReadFile(path)
			if err != nil {
				return err
			}
			for _, activity := range tdb.Acts.Act {
				for _, lap := range activity.Laps {
					database.laps = append(database.laps, lap)
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	database.total = database.Total(NullFilter)
	return nil
}
