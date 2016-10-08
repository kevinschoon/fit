package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gonum/matrix"
	mtx "github.com/gonum/matrix/mat64"
	"io"
	"net/url"
	"strconv"
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
	Datasets []struct {
		Name    string   // Name of the dataset
		Columns []string // Columns within the dataset to query
	}
	Function *Function // Aggregation Function
	fnStr    string    // Function name
	Max      int       // Maximum aggregation value
	Col      int       // Column to aggregate on
}

// Len returns the length of the Query
func (query Query) Len() int {
	return len(query.Datasets)
}

// Columns returns a flattened ordered
// array of Column names
func (query Query) Columns() []string {
	columns := make([]string, 0)
	for _, dataset := range query.Datasets {
		for _, column := range dataset.Columns {
			columns = append(columns, column)
		}
	}
	return columns
}

// QueryStr returns a valid URL query string
func (query Query) QueryStr() string {
	values := url.Values{}
	for _, dataset := range query.Datasets {
		args := make([]string, len(dataset.Columns)+1)
		args[0] = dataset.Name
		for i, column := range dataset.Columns {
			args[i+1] = column
		}
		values.Add("q", strings.TrimRight(strings.Join(args, ","), ","))
	}
	if query.Function != nil {
		values.Add("fn", query.fnStr)
		values.Add("max", fmt.Sprintf("%d", query.Max))
		values.Add("col", fmt.Sprintf("%d", query.Col))
	}
	return values.Encode()
}

// NewQuery constructs a Query from the provided
// args and optionally specified function.
// If function is specified the query returns
// aggregated
func NewQuery(args []string, function string, max, col int) *Query {
	query := &Query{
		Datasets: make([]struct {
			Name    string
			Columns []string
		}, len(args)),
		Max:   max,
		Col:   col,
		fnStr: function,
	}
	for i, arg := range args {
		split := strings.Split(arg, ",")
		if len(split) >= 1 {
			query.Datasets[i].Name = split[0]
		}
		if len(split) > 1 {
			query.Datasets[i].Columns = split[1:]
		}
	}
	// If a valid function is specified we apply aggregation
	// when the Query returns a result
	switch function {
	case "sum":
		query.Function = &Sum
	case "min":
		query.Function = &Min
	case "max":
		query.Function = &Max
	case "avg":
		query.Function = &Avg
	}
	return query
}

// NewQueryFromQS constructs a query from a url.URL
func NewQueryFromQS(u *url.URL) *Query {
	var args []string
	query := u.Query()
	if q, ok := query["q"]; ok {
		args = q
	}
	max, _ := strconv.ParseInt(query.Get("max"), 0, 64)
	col, _ := strconv.ParseInt(query.Get("col"), 0, 64)
	return NewQuery(args, query.Get("fn"), int(max), int(col))
}
