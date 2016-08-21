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

// Total represents aggregated tcx lap and track data
type Total struct {
	Start     time.Time
	TotalTime float64
	Dist      float64
}

type Totals []*Total

// Database is an singleton containing tcx lap data
type Database struct {
	laps tcx.Laps
}

// MetersTotal returns the total number of meters recorded
func (db *Database) Meters(fn Filter) int {
	var total float64
	for _, lap := range db.Laps(fn) {
		total += lap.Dist
	}
	return int(total)
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

// Aggr aggregates tcx laps by "year", "month", or "day"
func (db *Database) Aggr(fn Filter, bucket string) []*Total {
	t := make(map[int]*Total)
	for _, lap := range db.Laps(fn) {
		switch bucket {
		case "year": // TODO: Month, Day
			if total, ok := t[lap.Start.Year()]; ok {
				total.Dist += lap.Dist
				total.TotalTime += lap.TotalTime
			} else {
				t[lap.Start.Year()] = &Total{
					Start: time.Date(lap.Start.Year(),
						time.January, 0, 0, 0, 0, 0, time.UTC),
					Dist:      lap.Dist,
					TotalTime: lap.TotalTime,
				}
			}
		default:
			panic(fmt.Sprintf("Unknown bucket: %s", bucket))
		}
	}
	totals := make(Totals, 0)
	for _, value := range t {
		totals = append(totals, value)
	}
	fmt.Println(totals)
	return totals
}

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

// InitDB initializes the global "database" struct when provided
// a directory of files containing TCX data. It has been tested
// with TCX data provided by Google Checkout.
func InitDB(directory string) error {
	database = &Database{
		laps: make(tcx.Laps, 0),
	}
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
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
}
