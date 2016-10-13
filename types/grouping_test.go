package types

import (
	"fmt"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
	"time"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func TestGrouping(t *testing.T) {
	grouping := NewGrouping("Duration,0,2s")
	assert.Equal(t, 2*time.Second, grouping.Max)
	assert.Equal(t, 0, grouping.Index)
	assert.Equal(t, "Duration,0,2s", grouping.String())
	mx := mtx.NewDense(10, 10, nil)
	r, c := mx.Dims()
	start := time.Date(2001, time.January, 0, 0, 0, 1, 1, time.UTC)
	for i := 0; i < r; i++ {
		mx.Set(i, 0, float64(start.Unix()))
		mx.Set(i, 1, 1.0)
		start = start.Add(1 * time.Second)
		for j := 2; j < c; j++ {
			mx.Set(i, j, toFixed(rand.Float64(), 2))
		}
	}
	assert.Equal(t, mtx.Sum(mx.ColView(1)), 10.0)
	assert.True(t, mtx.Sum(mx.ColView(2)) < 10)
	views := grouping.Group(mx)
	assert.Len(t, views, 5)
	rows := 0
	for _, v := range views {
		fmt.Println(mtx.Formatted(v))
		r, _ := v.Dims()
		rows += r
	}
	assert.Equal(t, 10, r)
	//mx = apply(result, Avg)
	//r, c = mx.Dims()
	//assert.Equal(t, 5, r)
	//assert.Equal(t, 10, c)
	//fmt.Println(mtx.Formatted(mx))
}

func init() {
	rand.Seed(time.Now().Unix())
}
