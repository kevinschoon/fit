package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kevinschoon/gofit/models"
)

// Writer implements a function to write to a SQL database
type Writer interface {
	Write(*gorm.DB) error
}

// Reader implements functions to query and read from a SQL database
type Reader interface {
	Read(*gorm.DB, Query) (models.Series, error)
}

// New creates a new gorm DB and automigrates the specified objects
func New(path string, objs ...interface{}) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	db.DB().SetMaxIdleConns(0)
	db.DB().SetMaxOpenConns(0)
	if err = db.AutoMigrate(objs...).Error; err != nil {
		return nil, err
	}
	return db, nil
}

// Write executes the Writer function
func Write(db *gorm.DB, writer Writer) error {
	return writer.Write(db)
}

// Read executes the Query and Read functions
func Read(db *gorm.DB, query Query, reader Reader) (models.Series, error) {
	return reader.Read(db, query)
}
