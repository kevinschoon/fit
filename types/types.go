package types

import (
	"encoding/json"
	"errors"
	"github.com/gonum/matrix"
	mtx "github.com/gonum/matrix/mat64"
	"io"
	"net/url"
	"strings"
	"sync"
)

var (
	ErrAPI      = errors.New("api error")
	ErrNoData   = errors.New("no data")
	ErrNotFound = errors.New("not found")
	ErrBadQuery = errors.New("bad query")
)

type Client interface {
	Datasets() ([]*Dataset, error)
	Write(*Dataset) error
	Read(string) (*Dataset, error)
	Delete(string) error
	Query(Queries) (*Dataset, error)
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
	Mtx     []float64
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
		out.Mtx = ds.Mtx.RawMatrix().Data
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
			ds.Mtx = mtx.NewDense(ds.Stats.Rows, ds.Stats.Columns, in.Mtx)
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

// Query can be used to combine the results
// of multiple datasets into a single
// matrix of values
type Query struct {
	Name    string   // Dataset name
	Columns []string // Column names
}

type Queries []*Query

func (queries Queries) Columns() []string {
	columns := make([]string, 0)
	for _, query := range queries {
		for _, column := range query.Columns {
			columns = append(columns, column)
		}
	}
	return columns
}

func (queries Queries) Len() int {
	return len(queries)
}

// QueryStr returns a valid URL query string
func (queries Queries) QueryStr() string {
	values := url.Values{}
	for _, query := range queries {
		args := make([]string, len(query.Columns)+1)
		args[0] = query.Name
		for i := 0; i < len(query.Columns); i++ {
			args[i+1] = query.Columns[i]
		}
		values.Add("q", strings.TrimRight(strings.Join(args, ","), ","))
	}
	return values.Encode()
}

func NewQueries(args []string) Queries {
	queries := make(Queries, len(args))
	for i, arg := range args {
		split := strings.Split(arg, ",")
		if len(split) > 1 {
			queries[i] = &Query{
				Name:    split[0],
				Columns: split[1:],
			}
		}
	}
	return queries
}

func NewQueriesFromQS(u *url.URL) Queries {
	queries := Queries{}
	if q, ok := u.Query()["q"]; ok {
		queries = NewQueries(q)
	}
	return queries
}
