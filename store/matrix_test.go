package store

import (
	"fmt"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/kevinschoon/fit/loader"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadMatrix(t *testing.T) {
	reader := loader.NewMemReader(
		[]string{"V1", "V2"},
		[][]float64{
			[]float64{
				1.0,
				2.0,
			},
			[]float64{
				3.0,
				4.0,
			},
		})
	m, err := ReadMatrix(reader)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, m.At(0, 0))
	assert.Equal(t, 2.0, m.At(0, 1))
	assert.Equal(t, 3.0, m.At(1, 0))
	assert.Equal(t, 4.0, m.At(1, 1))
}

func TestMatrix(t *testing.T) {
	m := mtx.NewDense(2, 2, []float64{1.0, 1.0, 1.0, 1.0})
	assert.Equal(t, 4.0, mtx.Sum(m))
	assert.Equal(t, 2, m.ColView(1).Len())
	d := mtx.Col(nil, 0, m)
	assert.Equal(t, 2, len(d))
	a := mtx.NewDense(2, 4, []float64{
		1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
	})
	b := mtx.NewDense(2, 4, []float64{
		2.0, 2.0, 2.0, 2.0,
		2.0, 2.0, 2.0, 2.0,
	})
	fn := func(i, j int, v float64) float64 {
		fmt.Println(i, j, v)
		return 0.0
	}
	fmt.Println(a.Dims())
	fmt.Println(a.T().Dims())
	a.Apply(fn, b)
}
