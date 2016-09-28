package store

import (
	"errors"
	mtx "github.com/gonum/matrix/mat64"
	"io"
)

var ErrUnequalValues = errors.New("unequal value size")

func ReadMatrix(reader Reader) (*mtx.Dense, error) {
	var (
		values []float64
		r      int
		c      int
	)
	c = len(reader.Columns())
	for {
		v, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(v) != c {
			return nil, ErrUnequalValues
		}
		for _, value := range v {
			values = append(values, value)
		}
		r++
	}
	return mtx.NewDense(r, len(reader.Columns()), values), nil
}
