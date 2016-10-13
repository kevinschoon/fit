package types

import (
	mtx "github.com/gonum/matrix/mat64"
	"strings"
)

type Function struct {
	Name string
}

func (fn Function) apply(mx []mtx.Matrix, f func(mtx.Matrix) float64) *mtx.Dense {
	var (
		cols   int
		result *mtx.Dense
	)
	for i, view := range mx {
		if result == nil {
			_, cols = view.Dims()
			result = mtx.NewDense(len(mx), cols, nil)
		}
		// TODO: Determine if there is a better way to accomplish
		// this without type assertion
		other := view.(*mtx.Dense)
		for j := 0; j < cols; j++ {
			result.Set(i, j, f(other.ColView(j)))
		}
	}
	return result
}

func (fn Function) Apply(mx []mtx.Matrix) *mtx.Dense {
	switch strings.ToLower(fn.Name) {
	case "min":
		return fn.apply(mx, mtx.Min)
	case "max":
		return fn.apply(mx, mtx.Max)
	case "sum":
		return fn.apply(mx, mtx.Sum)
	default: // Use the average
		if len(mx) > 0 {
			return fn.apply(mx, func(o mtx.Matrix) float64 {
				r, _ := o.Dims()
				return mtx.Sum(o) / float64(r)
			})
		}
	}
	return mtx.NewDense(0, 0, nil)
}
