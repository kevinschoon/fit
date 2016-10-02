package chart

import (
	mtx "github.com/gonum/matrix/mat64"
	"github.com/gonum/plot/plotter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLines(t *testing.T) {
	mx := mtx.NewDense(3, 3, []float64{
		1.0, 1.0, 2.0,
		2.0, 3.0, 4.0,
		3.0, 5.0, 6.0,
	})
	data := GetLines(mx, []string{"x", "y", "z"})
	assert.Len(t, data, 4)
	assert.IsType(t, "", data[0])
	name := data[0].(string)
	assert.Equal(t, "y", name)
	assert.IsType(t, plotter.XYs{}, data[1])
	xys := data[1].(plotter.XYs)
	assert.Equal(t, 3, xys.Len())
	x, y := xys.XY(0)
	assert.Equal(t, 1.0, x)
	assert.Equal(t, 1.0, y)
	x, y = xys.XY(1)
	assert.Equal(t, 2.0, x)
	assert.Equal(t, 3.0, y)
	x, y = xys.XY(2)
	assert.Equal(t, 3.0, x)
	assert.Equal(t, 5.0, y)
	assert.IsType(t, "", data[2])
	name = data[2].(string)
	assert.Equal(t, "z", name)
	assert.IsType(t, plotter.XYs{}, data[3])
	xys = data[3].(plotter.XYs)
	assert.Equal(t, 3, xys.Len())
	x, y = xys.XY(0)
	assert.Equal(t, 1.0, x)
	assert.Equal(t, 2.0, y)
	x, y = xys.XY(1)
	assert.Equal(t, 2.0, x)
	assert.Equal(t, 4.0, y)
	x, y = xys.XY(2)
	assert.Equal(t, 3.0, x)
	assert.Equal(t, 6.0, y)
}
