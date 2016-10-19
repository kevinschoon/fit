package types

import (
	"encoding/json"
	"errors"
	"github.com/gonum/matrix"
	mtx "github.com/gonum/matrix/mat64"
	"io"
	"math"
	"sync"
)

var (
	ErrAPI      = errors.New("api error")
	ErrNoData   = errors.New("no data")
	ErrNotFound = errors.New("not found")
	ErrBadQuery = errors.New("bad query")
)

// Work around for handling NaN values in JSON
// https://github.com/golang/go/issues/3480
type value float64

func (v value) MarshalJSON() ([]byte, error) {
	if math.IsNaN(float64(v)) {
		return json.Marshal(nil)
	}
	return json.Marshal(float64(v))
}

func (v *value) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, nil); err != nil {
		val := float64(0.0)
		if err := json.Unmarshal(data, &val); err == nil {
			*v = value(val)
		}
	} else {
		*v = value(math.NaN())
	}
	return nil
}

type values []value

type Client interface {
	Datasets() ([]*Dataset, error)
	Write(*Dataset) error
	Delete(string) error
	Query(*Query) (*Dataset, error)
}

// Stats contain statistics about the
// underlying data in a dataset
type Stats struct {
	Rows    int
	Columns int
}

type dataset struct {
	Name    string
	Columns []string
	Stats   *Stats
	Mtx     []value
}

// Dataset consists of a name and
// an ordered array of column names
type Dataset struct {
	Name       string     // Name of this dataset
	Columns    []string   // Ordered array of cols
	Mtx        *mtx.Dense `json:"-"` // Dense Matrix contains all values in the dataset
	Stats      *Stats
	lock       sync.RWMutex
	index      int
	WithValues bool
}

func (ds *Dataset) MarshalJSON() ([]byte, error) {
	ds.stats()
	out := &dataset{
		Name:    ds.Name,
		Columns: ds.Columns,
		Stats:   ds.Stats,
	}
	if ds.WithValues && ds.Mtx != nil {
		r, c := ds.Mtx.Dims()
		out.Mtx = make([]value, r*c)
		for i, val := range ds.Mtx.RawMatrix().Data {
			out.Mtx[i] = value(val)
		}
	}
	return json.Marshal(out)
}

func (ds *Dataset) UnmarshalJSON(data []byte) error {
	in := &dataset{}
	if err := json.Unmarshal(data, in); err != nil {
		return err
	}
	ds.Name = in.Name
	ds.Columns = in.Columns
	ds.Stats = in.Stats
	return matrix.Maybe(func() {
		if ds.WithValues && in.Mtx != nil {
			ds.Mtx = mtx.NewDense(ds.Stats.Rows, ds.Stats.Columns, nil)
			values := make([]float64, ds.Stats.Rows*ds.Stats.Columns)
			for i := 0; i < len(in.Mtx); i++ {
				values[i] = float64(in.Mtx[i])
			}
		}
	})
}

// stats updates the Stats struct
func (ds *Dataset) stats() {
	if ds.Stats == nil {
		ds.Stats = &Stats{}
	}
	if ds.Mtx != nil {
		ds.Stats.Rows, ds.Stats.Columns = ds.Mtx.Dims()
	}
}

// Len returns the length (number of rows) of the dataset
func (ds Dataset) Len() int {
	len := 0
	if ds.Mtx != nil {
		len, _ = ds.Mtx.Dims()
	}
	return len
}

// CPos returns the position of a column
// name in a dataset. If the column
// does not exist it returns -1
func (ds Dataset) CPos(name string) int {
	for i, col := range ds.Columns {
		if name == col {
			return i
		}
	}
	return -1
}

// Next returns the next row of values
// If all values have been traversed
// it returns io.EOF. Implements the
// loader.Reader interface
func (ds *Dataset) Next() ([]float64, error) {
	ds.lock.Lock()
	defer ds.lock.Unlock()
	if ds.Mtx == nil {
		return nil, ErrNoData
	}
	r, _ := ds.Mtx.Dims()
	if ds.index >= r {
		ds.index = 0
		return nil, io.EOF
	}
	rows := ds.Mtx.RawRowView(ds.index)
	ds.index++
	return rows, nil
}
