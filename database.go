package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kevinschoon/tcx"
	"time"
)

const (
	MaxItems int    = 100
	dbTime   string = "2006-01-02"  // Date format for SQL query
	qTime    string = "2006-Jan-02" // Date format URL query
)

var options *Options

type Options struct {
	DBPath *string
	Debug  *bool
}

func GetDB() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", *options.DBPath)
	if err != nil {
		return nil, err
	}
	db.LogMode(*options.Debug)
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(0)
	return db, nil
}

func InitDB() {
	db, err := GetDB()
	FailOnErr(err)
	defer db.Close()
	FailOnErr(db.AutoMigrate(
		&tcx.Activity{},
		&tcx.Lap{},
		&tcx.Track{},
		&tcx.Trackpoint{}).Error)
}

func BulkUpsert(db *gorm.DB, activities []tcx.Activity) (err error) {
	for _, activity := range activities {
		db.Create(&activity)
		if db.Error != nil {
			return db.Error
		}
	}
	return nil
}

func Activity(id uint, db *gorm.DB) (activity tcx.Activity) {
	activity.ID = id
	db.Preload("Laps").Find(&activity)
	if db.Error == nil {
		for i, lap := range activity.Laps {
			activity.Laps[i].Trk = &tcx.Track{}
			db.Where("lap_id = ?", lap.ID).Preload("Pt").Find(activity.Laps[i].Trk)
		}
	}
	return activity
}

func Activities(db *gorm.DB, fn func(*gorm.DB) *gorm.DB) (activities tcx.Acts, err error) {
	var (
		count int
		last  int
	)
	fn(db).Model(&tcx.Activity{}).Count(&count)
	for len(activities) < count { // TODO: Cleanup
		if db.Error != nil {
			return nil, db.Error
		}
		var results tcx.Acts
		fn(db).Limit(MaxItems).Offset(last).Preload("Laps.Trk").Find(&results)
		for _, result := range results {
			activities = append(activities, result)
		}
		last += MaxItems
	}
	return activities, err
}

func TrackPoints(db *gorm.DB, fn func(*gorm.DB) *gorm.DB) (pts tcx.Trackpoints) {
	fn(db).Find(&pts)
	return pts
}

func Between(start, end time.Time, column string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s >= ? AND %s <= ?", column, column), start.Format(dbTime), end.Format(dbTime))
	}
}

/*
func (t Total) Km() float64 {
	return t.Dist * 0.001
}

// Totals is an array of Total
type Totals []*Total

func (t Totals) Predict() stats.Series {
	series := make(stats.Series, len(t))
	for i := 0; i < len(t); i++ {
		series[i] = stats.Coordinate{t[i].TotalTime, t[i].Dist}
	}
	regressions, err := stats.LinearRegression(series)
	if err != nil {
		panic(err)
	}
	return regressions
}
*/
