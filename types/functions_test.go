package types

import (
	mtx "github.com/gonum/matrix/mat64"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func data(r, c int, rnd bool) []float64 {
	data := make([]float64, r*c)
	for i := 0; i < r*c; i++ {
		if rnd {
			data[i] = rand.Float64()
		} else {
			data[i] = 1.0
		}
	}
	return data
}

func TestAggregate(t *testing.T) {
	assert.Equal(t, 100.0, mtx.Sum(mtx.NewDense(10, 10, data(10, 10, false))))
	assert.Equal(t, 100.0, mtx.Sum(Aggregate(3, 0, Sum, mtx.NewDense(10, 10, data(10, 10, false)))))
	mx := Aggregate(3, 0, Function(mtx.Sum), mtx.NewDense(10, 10, data(10, 10, false)))
	r, c := mx.Dims()
	assert.Equal(t, 4, r)
	assert.Equal(t, 10, c)
	assert.Equal(t, 3.0, mx.At(0, 0))
	assert.Equal(t, 3.0, mx.At(1, 0))
	assert.Equal(t, 3.0, mx.At(2, 0))
	assert.Equal(t, 543910.0, mtx.Sum(Aggregate(7, 0, Sum, mtx.NewDense(54391, 10, data(54391, 10, false)))))
}

func BenchmarkAggregate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Aggregate(5, 0, Sum, mtx.NewDense(5000, 10, data(5000, 10, false)))
	}
}

func BenchmarkAggregateRnd(b *testing.B) {
	d := data(5000, 10, true)
	for i := 0; i < b.N; i++ {
		Aggregate(5, 0, Sum, mtx.NewDense(5000, 10, d))
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}
