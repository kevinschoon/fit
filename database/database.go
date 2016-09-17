package database

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/kevinschoon/gofit/models"
	"time"
)

var ErrSeriesNotFound = errors.New("series not found")
var ErrSeriesNoData = errors.New("series has no data")

// DB is a wrapper for persisting Series data into Bolt
// Series are stored in two buckets:
// models: Contains the Series Keys and other information
// $NAME_data: Contains all of the Values in the series object
// Internally the values are broken up into second aggregated keys
// Currently updates are *NOT* supported
type DB struct {
	bolt  *bolt.DB
	debug bool
}

// Series returns all Series in the database without their Values
func (db *DB) Series() (series []*models.Series, err error) {
	err = db.bolt.View(func(tx *bolt.Tx) error {
		// If b == nil no Series exist yet
		if b := tx.Bucket([]byte("models")); b != nil {
			cursor := b.Cursor()
			for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
				s := &models.Series{}
				if err := json.Unmarshal(v, s); err != nil {
					return err
				}
				series = append(series, s)
			}
		}
		return nil
	})
	return series, err
}

// WriteSeries persists series to the database
func (db *DB) WriteSeries(series []*models.Series) error {
	// Resize the Series array into second-level aggregations
	series = models.Resize(series, 1*time.Second)
	// Get the name of the series
	name := series[0].Name
	return db.bolt.Update(func(tx *bolt.Tx) error {
		// Create the models bucket if it doesn't exist yet
		b, _ := tx.CreateBucketIfNotExists([]byte("models"))
		// Marshal the series without it's values to JSON
		raw, err := json.Marshal(series[0])
		if err != nil {
			return err
		}
		// Series model is stored in a key with the same name
		if err := b.Put([]byte(name), raw); err != nil {
			return err
		}
		// Get the bucket for storing Series data
		b, err = tx.CreateBucketIfNotExists([]byte(fmt.Sprintf("%s_data", name)))
		// Possible error if the key is invalid
		if err != nil {
			return err
		}
		for _, s := range series {
			// Dump all of the Values in the series and Marshal to JSON
			// TODO: Serialize to something faster/more efficent than JSON
			raw, err := json.Marshal(s.Dump())
			if err != nil {
				return err
			}
			// Values for each Series are written to a key
			// with a timestamp corresponding to the time
			// of the first set of Values for this Series
			if err := b.Put([]byte(s.Start().UTC().Format(time.RFC3339)), raw); err != nil {
				return err
			}
		}
		return nil
	})
}

// ReadSeries reads a Series from the database
func (db *DB) ReadSeries(name string, start, end time.Time) (series []*models.Series, err error) {
	err = db.bolt.View(func(tx *bolt.Tx) error {
		// Get the models bucket
		b := tx.Bucket([]byte("models"))
		if b == nil {
			// No Series have been saved
			return ErrSeriesNotFound
		}
		raw := b.Get([]byte(name))
		if raw == nil {
			// Series with the provided name was not found
			return ErrSeriesNotFound
		}
		// Unmarshal the series without data
		s := &models.Series{}
		if err = json.Unmarshal(raw, s); err != nil {
			return err
		}
		// Possible error if key is not valid
		b = tx.Bucket([]byte(fmt.Sprintf("%s_data", s.Name)))
		if b == nil {
			// Series has no data
			return ErrSeriesNoData
		}
		cursor := b.Cursor()
		// Perform a Range scan against the bucket
		// https://github.com/boltdb/bolt#range-scans
		min, max := []byte(start.Format(time.RFC3339)), []byte(end.Format(time.RFC3339))
		for k, v := cursor.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = cursor.Next() {
			// Shallow copy loaded Series model
			s = models.Copy(s)
			// Each set of values is put in a seperate Series
			// Series can be Resized together by the caller
			// Create an empty [][]Value container
			values := models.Values{}
			if err := json.Unmarshal(v, &values); err != nil {
				return err
			}
			// Import all the loaded values to the Series
			s.Import(values)
			// Append this series to the array
			series = append(series, s)
		}
		return nil
	})
	return series, err
}

func (db *DB) Close() {
	db.bolt.Close()
}

// New returns a new DB object for accesing
// Series data. It provides a wrapper around BoltDB
func New(path string, debug bool) (*DB, error) {
	b, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	db := &DB{
		bolt:  b,
		debug: debug,
	}
	return db, nil
}
