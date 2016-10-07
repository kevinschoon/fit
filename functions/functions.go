package functions

import (
	mtx "github.com/gonum/matrix/mat64"
)

type Function func(mtx.Matrix) float64

var (
	Sum = Function(mtx.Sum)
	Min = Function(mtx.Min)
	Max = Function(mtx.Max)
	Avg = Function(func(mx mtx.Matrix) float64 {
		len, _ := mx.Dims()
		return Sum(mx) / float64(len)
	})
)

// Aggregate slides a view across the provided
// Matrix tracking the sum of col.
// When col total exceeds max apply fn
// to all the values in the current view.
func Aggregate(max int, col int, fn Function, other *mtx.Dense) *mtx.Dense {
	rows, cols := other.Dims()
	if float64(max) > mtx.Sum(other.ColView(col).T()) {
		return other
	}
	result := mtx.NewDense(0, cols, nil)
	var view mtx.Matrix
	var totl float64
	for i, j := 0, 1; i+j <= rows; j++ {
		view = other.View(i, 0, j, cols)
		if totl+view.At(j-1, col) >= float64(max) {
			result = result.Grow(1, 0).(*mtx.Dense)
			result.SetRow(i/j, apply(view, fn))
			totl = float64(0)
			i, j = i+j, 0
		} else {
			totl += view.At(j-1, col)
		}
		if i+j >= rows { // Final iteration
			result = result.Grow(1, 0).(*mtx.Dense)
			r, _ := result.Dims()
			result.SetRow(r-1, apply(view, fn))
		}
	}
	return result
}

func apply(mx mtx.Matrix, fn Function) []float64 {
	_, c := mx.Dims()
	result := make([]float64, c)
	m := mx.(*mtx.Dense)
	for i := 0; i < c; i++ {
		result[i] = fn(m.ColView(i))
	}
	return result
}
