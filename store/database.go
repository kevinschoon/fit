package store

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	mtx "github.com/gonum/matrix/mat64"
	"strings"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
	ErrBadQuery = errors.New("bad query")
)

// Query can be used to combine the results
// of multiple datasets into a single
// matrix of values
type Query struct {
	Fn      string   // TODO
	Name    string   // Dataset name
	Columns []string // Column names
}

func QueryFromArgs(args []string) []*Query {
	q := make([]*Query, len(args))
	for i, arg := range args {
		split := strings.Split(arg, ",")
		q[i] = &Query{
			Name:    split[0],
			Columns: split[1:],
		}
	}
	return q
}

type DB struct {
	bolt *bolt.DB
}

func (db *DB) Datasets() (datasets []*Dataset, err error) {
	err = db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("datasets"))
		if b == nil { // No datasets have been saved
			return nil
		}
		cursor := b.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			ds := &Dataset{}
			if err := json.Unmarshal(v, ds); err != nil {
				return err
			}
			datasets = append(datasets, ds)
		}
		return nil
	})
	return datasets, err
}

func (db *DB) Write(ds *Dataset) (err error) {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("datasets"))
		if err != nil {
			return err
		}
		raw, err := json.Marshal(ds)
		if err != nil {
			return err
		}
		if err = b.Put([]byte(ds.Name), raw); err != nil {
			return err
		}
		b, err = tx.CreateBucketIfNotExists([]byte("matricies"))
		if err != nil {
			return err
		}
		raw, err = ds.Mtx.MarshalBinary()
		if err != nil {
			return err
		}
		return b.Put([]byte(ds.Name), raw)
	})
}

func (db *DB) Read(name string) (ds *Dataset, err error) {
	err = db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("datasets"))
		if b == nil { // No datasets have been saved
			return nil
		}
		raw := b.Get([]byte(name))
		if raw == nil {
			return ErrNotFound
		}
		ds = &Dataset{}
		if err = json.Unmarshal(raw, ds); err != nil {
			return err
		}
		b = tx.Bucket([]byte("matricies"))
		if b == nil {
			return ErrNotFound
		}
		raw = b.Get([]byte(name))
		if raw == nil {
			return ErrNotFound
		}
		ds.Mtx = mtx.NewDense(0, 0, nil)
		return ds.Mtx.UnmarshalBinary(raw)
	})
	return ds, err
}

func (db *DB) Query(queries ...*Query) (*Dataset, error) {
	var (
		rows int // Row count for new dataset
		cols int // Col count for new dataset
	)
	// The new resulting dataset
	ds := &Dataset{
		Name:    "QueryResult",
		Columns: make([]string, 0),
	}
	// Empty array of Vectors where each
	// is a column from the queries
	vectors := make([]*mtx.Vector, 0)
	// Range each query
	for _, query := range queries {
		// Query for the other dataset
		other, err := db.Read(query.Name)
		if err != nil {
			return nil, err
		}
		// Resulting matrix should have the sum of
		// the number of rows from each matrix that
		// is queried
		r, _ := other.Mtx.Dims()
		rows += r
		// Range each column in the query
		for _, name := range query.Columns {
			// Get the position (index) of the column
			pos := other.CPos(name)
			// If the returned position is a negative
			// number the column does not exist
			if pos < 0 {
				return nil, ErrNotFound
			}
			// Append the column to vectors array
			vectors = append(vectors, other.Mtx.ColView(pos))
			// Add the column name to the resulting dataset
			ds.Columns = append(ds.Columns, name)
		}
	}
	// Resulting number of columns is equal to
	// the amount that were queried for
	cols = len(vectors)
	// Create a new matrix zeroed Matrix
	ds.Mtx = mtx.NewDense(rows, cols, nil)
	// Fill the matrix with values from each column vector
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if vectors[j].Len() > i {
				ds.Mtx.Set(i, j, vectors[j].At(i, 0))
			} // Zeros are left for missing data
		}
	}
	return ds, nil
}

func (db *DB) Close() {
	db.bolt.Close()
}

// New returns a new DB object for accesing
// Series data. It provides a wrapper around BoltDB
func NewDB(path string) (*DB, error) {
	b, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	db := &DB{
		bolt: b,
	}
	return db, nil
}
