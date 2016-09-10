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

var ErrCollectionNotFound = errors.New("collection not found")

type DB struct {
	b *bolt.DB
}

// Write executes the Writer function
func (db *DB) Write(name string, collection *models.Collection) error {
	return db.b.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}
		for _, series := range collection.Series {
			raw, err := json.Marshal(series.Values)
			if err != nil {
				return err
			}
			if err := b.Put([]byte(series.Time.Format(time.RFC3339)), raw); err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *DB) Read(name string, start, end time.Time) (*models.Collection, error) {
	var cursor *bolt.Cursor
	collection := &models.Collection{
		Name: name,
	}
	err := db.b.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(name)); b != nil {
			cursor = b.Cursor()
		} else {
			return ErrCollectionNotFound
		}
		min, max := []byte(start.Format(time.RFC3339)), []byte(end.Format(time.RFC3339))
		for k, v := cursor.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = cursor.Next() {
			fmt.Printf("%s\n", k)
			data := [][]models.Value{}
			if err := json.Unmarshal(v, &data); err != nil {
				return err
			}
			startTime, err := time.Parse(time.RFC3339, string(k))
			if err != nil {
				return err
			}
			for _, values := range data {
				collection.Add(startTime, values)
			}
		}
		return nil
	})
	return collection, err
}

func (db *DB) Close() {
	db.b.Close()
}

// New returns a new DB object
func New(path string) (*DB, error) {
	b, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	db := &DB{
		b: b,
	}
	return db, nil
}
