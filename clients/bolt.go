package clients

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/kevinschoon/fit/types"
	"time"
)

// BoltClient implements the types.Client
// interface with a BoltDB backend.
// BoltClient is the primary database
// backend for Fit. It can also be used
// directly from the command line.
type BoltClient struct {
	bolt *bolt.DB
}

var (
	dsBucket = []byte("datasets")
	mxBucket = []byte("matricies")
)

func (c *BoltClient) Datasets() (datasets []*types.Dataset, err error) {
	err = c.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dsBucket)
		cursor := b.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			ds := &types.Dataset{}
			if err := json.Unmarshal(v, ds); err != nil {
				return err
			}
			datasets = append(datasets, ds)
		}
		return nil
	})
	return datasets, err
}

func (c *BoltClient) Write(ds *types.Dataset) (err error) {
	return c.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(dsBucket)
		raw, err := json.Marshal(ds)
		if err != nil {
			return err
		}
		if err = b.Put([]byte(ds.Name), raw); err != nil {
			return err
		}
		if ds.Mtx == nil { // No matricies attached to this dataset
			return nil
		}
		b = tx.Bucket(mxBucket)
		raw, err = ds.Mtx.MarshalBinary()
		if err != nil {
			return err
		}
		return b.Put([]byte(ds.Name), raw)
	})
}

func (c *BoltClient) read(name string) (ds *types.Dataset, err error) {
	if err = c.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dsBucket)
		raw := b.Get([]byte(name))
		if raw == nil {
			return types.ErrNotFound
		}
		ds = &types.Dataset{}
		if err = json.Unmarshal(raw, ds); err != nil {
			return err
		}
		b = tx.Bucket(mxBucket)
		raw = b.Get([]byte(name))
		if raw == nil {
			return nil // No matricies attached to the dataset
		}
		ds.Mtx = mtx.NewDense(0, 0, nil)
		return ds.Mtx.UnmarshalBinary(raw)
	}); err != nil {
		return nil, err
	}
	return ds, nil
}

func (c *BoltClient) Delete(name string) error {
	return c.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(dsBucket)
		if err := b.Delete([]byte(name)); err != nil {
			return err
		}
		b = tx.Bucket(mxBucket)
		return b.Delete([]byte(name))
	})
}

// Query finds all of the datasets contained
// in Queries and returns a combined dataset
// for each column in the search. The values
// from each dataset are stored entirely in
// memory until the query is complete. This
// means that the total size of all datasets
// queried cannot exceed the total system
// memory. The resulting dataset columns
// will be ordered in the same order they
// were queried for.
func (c *BoltClient) Query(query *types.Query) (*types.Dataset, error) {
	var (
		rows  int            // Row count for new dataset
		cols  int            // Col count for new dataset
		other *types.Dataset // Dataset currently being processed
	)
	// The new resulting dataset
	ds := &types.Dataset{
		Name:    "QueryResult",
		Columns: make([]string, 0),
	}
	// Empty array of Vectors where each
	// is a column from the queries
	vectors := make([]*mtx.Vector, 0)
	// Map of datasets already processed
	processed := make(map[string]*types.Dataset)
	// Range each dataset in the query
	for _, dataset := range query.Datasets {
		columns := dataset.Columns
		//columns := query.Columns(name)
		// Check to see if a query for this dataset
		// has already been executed
		if _, ok := processed[dataset.Name]; !ok {
			// Query for the other dataset
			other, err := c.read(dataset.Name)
			if err != nil {
				return nil, err
			}
			// Resulting matrix should have the sum of
			// the number of rows from each unique
			// dataset matrix that is queried
			r, _ := other.Mtx.Dims()
			rows += r
			// Add this dataset to the map
			// so it is not queried again
			processed[dataset.Name] = other
		}
		// The other dataset we are querying
		other = processed[dataset.Name]
		// If this is a wild card search
		// set columns to equal all available
		// columns in the dataset
		if len(columns) == 1 {
			if columns[0] == "*" {
				columns = other.Columns
			}
		}
		// Range each column in the query
		for _, name := range columns {
			// Get the position (index) of the column
			pos := other.CPos(name)
			// If the returned position is a negative
			// number the column does not exist
			if pos < 0 {
				return nil, types.ErrNotFound
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
	// Apply any other query options to the resulting dataset
	ds.Mtx = query.Apply(ds.Mtx)
	return ds, nil
}

func (c *BoltClient) Close() {
	c.bolt.Close()
}

func NewBoltClient(path string) (types.Client, error) {
	b, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	c := &BoltClient{
		bolt: b,
	}
	// Initialize buckets
	err = c.bolt.Update(func(tx *bolt.Tx) error {
		if _, err = tx.CreateBucketIfNotExists(dsBucket); err != nil {
			return err
		}
		if _, err = tx.CreateBucketIfNotExists(mxBucket); err != nil {
			return err
		}
		return nil
	})
	return types.Client(c), err
}
