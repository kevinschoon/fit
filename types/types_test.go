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

func TestNewQuery(t *testing.T) {
	args := []string{"D0,x,y", "D1,x,z", "D1,x"}
	query := NewQuery(args, "sum", 10, 1)
	assert.Equal(t, 10, query.Max)
	assert.Equal(t, 1, query.Col)
	assert.Exactly(t, query.Function, &Sum)
	assert.Equal(t, 3, query.Len())
	assert.Equal(t, "D0", query.Datasets[0].Name)
	assert.Len(t, query.Datasets[0].Columns, 2)
	assert.Equal(t, "D1", query.Datasets[1].Name)
	assert.Len(t, query.Datasets[1].Columns, 2)
	assert.Equal(t, "D1", query.Datasets[1].Name)
	assert.Len(t, query.Datasets[2].Columns, 1)
	assert.Equal(t, "col=1&fn=sum&max=10&q=D0%2Cx%2Cy&q=D1%2Cx%2Cz&q=D1%2Cx", query.QueryStr())
}

func TestNewQueryFromQS(t *testing.T) {
	u, err := url.Parse("http://localhost?q=Fuu,x&q=Bar,y,z&fn=sum&max=10&col=1")
	assert.NoError(t, err)
	query := NewQueryFromQS(u)
	assert.Exactly(t, query.Function, &Sum)
	assert.Equal(t, 2, query.Len())
	assert.Equal(t, 1, query.Col)
	assert.Equal(t, 10, query.Max)
	assert.Equal(t, "Fuu", query.Datasets[0].Name)
	assert.Len(t, query.Datasets[0].Columns, 1)
	assert.Equal(t, "x", query.Datasets[0].Columns[0])
	assert.Equal(t, "Bar", query.Datasets[1].Name)
	assert.Len(t, query.Datasets[1].Columns, 2)
	assert.Equal(t, "y", query.Datasets[1].Columns[0])
	assert.Equal(t, "z", query.Datasets[1].Columns[1])
}
