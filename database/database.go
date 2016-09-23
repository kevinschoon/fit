package database

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/kevinschoon/fit/models"
	"log"
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
	if err := db.bolt.Update(func(tx *bolt.Tx) error {
		// Create the models bucket if it doesn't exist yet
		b, _ := tx.CreateBucketIfNotExists([]byte("models"))
		// Marshal the series without it's values to JSON
		raw, err := json.Marshal(series[0])
		if err != nil {
			return err
		}
		// Series model is stored in a key with the same name
		return b.Put([]byte(name), raw)
	}); err != nil {
		return err
	}
	// Persist all values in the series
	// Series values are stored in a bucket $SERIES_data
	bucket := fmt.Sprintf("%s_data", name)
	fmt.Println(bucket)
	for _, s := range series {
		if err := db.WriteValues(bucket, s.Start(), s.Dump()); err != nil {
			return err
		}
	}
	return nil
}

// WriteValues persists series Values to the database
func (db *DB) WriteValues(name string, start time.Time, values models.Values) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		// Create a new values bucket if it does not exist
		b, err := tx.CreateBucketIfNotExists([]byte(name))
		// Possible error if the key is invalid
		if err != nil {
			return err
		}
		// Marshal all Values to JSON
		// TODO: Serialize to something faster/more efficent than JSON
		raw, err := json.Marshal(values)
		if err != nil {
			return err
		}
		// Values are written to a key with a timestamp corresponding
		// to the start time provided. When called by the WriteSeries
		// function this will be the start time of the Series.
		return b.Put([]byte(start.Format(time.RFC3339)), raw)
	})
}

// ReadSeries reads a Series from the database
// ALL values must be capable of fitting in memory
func (db *DB) ReadSeries(name string, start, end time.Time) (series []*models.Series, err error) {
	log.Printf("READ [%s] - (%s-%s)", name, start.String(), end.String())
	template := &models.Series{}
	// Attempt to read a series object with the provided name
	if err := db.bolt.View(func(tx *bolt.Tx) error {
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
		// Unmarshal the series without data to the "template" Series
		return json.Unmarshal(raw, template)
	}); err != nil {
		return series, err
	}
	// Create an unbuffered channel
	ch := make(chan models.Values, 0)
	// ReadValues in a seperate routine
	go func() {
		// Set err to the result of ReadValues
		err = db.ReadValues(fmt.Sprintf("%s_data", name), start, end, ch)
	}()
	// Block until all values have been read
	for {
		if values, ok := <-ch; ok {
			// Create a new series for each set of values
			s := models.Copy(template)
			// Import the values into the series
			s.Import(values)
			// Add the series to array
			series = append(series, s)
		} else { // Channel was closed
			break
		}
	}
	return series, err
}

// ReadValues returns values within the requested range from the database
func (db *DB) ReadValues(name string, start, end time.Time, ch chan models.Values) error {
	defer close(ch) // Close the channel if there is an error or returned normally
	return db.bolt.View(func(tx *bolt.Tx) error {
		// Possible error if key is not valid
		b := tx.Bucket([]byte(name))
		if b == nil {
			// Series has no data
			return ErrSeriesNoData
		}
		cursor := b.Cursor()
		// Perform a Range scan against the bucket
		// https://github.com/boltdb/bolt#range-scans
		min, max := []byte(start.Format(time.RFC3339)), []byte(end.Format(time.RFC3339))
		for k, v := cursor.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = cursor.Next() {
			// Each set of values is put in a seperate Series
			// Series can be Resized together by the caller
			// Create an empty [][]Value container
			values := models.Values{}
			if err := json.Unmarshal(v, &values); err != nil {
				return err
			}
			// Send the values back
			ch <- values
		}
		return nil
	})
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
