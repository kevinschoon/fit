package store

import mtx "github.com/gonum/matrix/mat64"

// Dataset consists of a name and
// an ordered array of column names
type Dataset struct {
	Name    string     // Name of this dataset
	Columns []string   // Ordered array of cols
	Mtx     *mtx.Dense `json:"-"` // Dense Matrix contains all values in the dataset
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

// Reader provides an iterative interface
// for loading sets of float64 values
type Reader interface {
	Next() ([]float64, error)
	Columns() []string
	Close() error
}
