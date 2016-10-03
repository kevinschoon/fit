package types

import (
	"encoding/json"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestDatasetJSON(t *testing.T) {
	ds := &Dataset{
		Name:       "TestDataset",
		Columns:    []string{"V1", "V2"},
		Mtx:        mtx.NewDense(2, 2, []float64{1.0, 1.0, 2.0, 2.0}),
		WithValues: true,
	}
	raw, err := json.Marshal(ds)
	assert.NoError(t, err)
	out := &Dataset{WithValues: true}
	assert.NoError(t, json.Unmarshal(raw, out))
	assert.Equal(t, ds.Name, out.Name)
	assert.Equal(t, ds.Columns[0], out.Columns[0])
	assert.Equal(t, ds.Columns[1], out.Columns[1])
	assert.Equal(t, ds.Mtx.At(1, 1), out.Mtx.At(1, 1))
}

func TestQueries(t *testing.T) {
	args := []string{"D0,x,y", "D1,x,z"}
	queries := NewQueries(args)
	assert.Equal(t, 2, len(queries))
	assert.Equal(t, "D0", queries[0].Name)
	assert.Equal(t, "D1", queries[1].Name)
	columns := queries.Columns()
	assert.Equal(t, 4, len(columns))
	assert.Equal(t, "x", columns[0])
	assert.Equal(t, "y", columns[1])
	assert.Equal(t, "x", columns[2])
	assert.Equal(t, "z", columns[3])
	assert.Equal(t, "q=D0%2Cx%2Cy&q=D1%2Cx%2Cz", queries.QueryStr())
}

func TestQueriesFromQS(t *testing.T) {
	u, err := url.Parse("http://localhost?q=Fuu,x&q=Bar,y,z")
	assert.NoError(t, err)
	queries := NewQueriesFromQS(u)
	assert.Equal(t, 2, queries.Len())
	columns := queries.Columns()
	assert.Len(t, columns, 3)
	assert.Equal(t, "x", columns[0])
	assert.Equal(t, "y", columns[1])
	assert.Equal(t, "z", columns[2])
}
