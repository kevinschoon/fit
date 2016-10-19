package types

import (
	"encoding/json"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestDatasetJSON(t *testing.T) {
	ds := &Dataset{
		Name:       "TestDataset",
		Columns:    []string{"V1", "V2"},
		Mtx:        mtx.NewDense(3, 2, []float64{1.0, 1.0, 2.0, 2.0, 3.0, math.NaN()}),
		WithValues: true,
	}
	raw, err := json.Marshal(ds)
	assert.NoError(t, err)
	out := &Dataset{WithValues: true}
	assert.NoError(t, json.Unmarshal(raw, out))
	assert.Equal(t, ds.Name, out.Name)
	assert.Equal(t, ds.Columns[0], out.Columns[0])
	assert.Equal(t, ds.Columns[1], out.Columns[1])
	assert.Equal(t, 1.0, ds.Mtx.At(0, 0))
	assert.Equal(t, 1.0, ds.Mtx.At(0, 1))

	assert.Equal(t, 2.0, ds.Mtx.At(1, 0))
	assert.Equal(t, 2.0, ds.Mtx.At(1, 1))

	assert.Equal(t, 3.0, ds.Mtx.At(2, 0))
	assert.True(t, math.IsNaN(ds.Mtx.At(2, 1)))
}
