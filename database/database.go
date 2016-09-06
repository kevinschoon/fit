package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kevinschoon/gofit/models"
	"sort"
)

// Writer implements a function to write to a SQL database
type Writer interface {
	Write(*gorm.DB) error
}

// Reader implements functions to query and read from a SQL database
type Reader interface {
	Read(*gorm.DB, models.Query) (models.Serieser, error)
}

// Migrater returns types for Gorm automigration
type Migrater interface {
	Types() []interface{}
}

// New creates a new gorm DB and automigrates the specified objects
func New(path string, migrater Migrater) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(0)
	if err = db.AutoMigrate(migrater.Types()...).Error; err != nil {
		return nil, err
	}
	return db, nil
}

// Write executes the Writer function
func Write(db *gorm.DB, writer Writer) error {
	return writer.Write(db)
}

// Read executes the Query and Read functions
func Read(db *gorm.DB, query models.Query, reader Reader) (*models.Series, error) {
	data, err := reader.Read(db, query)
	if err != nil {
		return nil, err
	}
	series := data.Series()
	series.Rows = series.Rows.RollUp(query.Precision)
	sort.Sort(sort.Reverse(series.Rows))
	if query.Order == "reverse" {
		sort.Reverse(series.Rows)
	}
	return series, nil
}
