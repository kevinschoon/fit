package clients

import (
	mtx "github.com/gonum/matrix/mat64"
	"github.com/kevinschoon/fit/types"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

var cleanup bool = false

func NewTestMatrix(r, c int) *mtx.Dense {
	values := make([]float64, r*c)
	for i := range values {
		values[i] = rand.NormFloat64()
	}
	return mtx.NewDense(r, c, values)
}

func NewTestDB(t *testing.T) (types.Client, func()) {
	f, err := ioutil.TempFile("/tmp", "fit-test")
	if err != nil {
		t.Error(err)
	}
	db, err := NewBoltClient(f.Name())
	assert.NoError(t, err)
	return db, func() {
		if cleanup {
			os.Remove(f.Name())
		}
	}
}

func TestDatasets(t *testing.T) {
	db, cleanup := NewTestDB(t)
	defer cleanup()
	datasets, err := db.Datasets()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(datasets))
}

func TestReadWrite(t *testing.T) {
	db, cleanup := NewTestDB(t)
	defer cleanup()
	dsA := &types.Dataset{
		Name: "TestReadWrite",
		Columns: []string{
			"V1", "V2", "V3", "V4",
			"V5", "V6", "V7", "V8",
		},
		Mtx: NewTestMatrix(128, 8),
	}
	assert.NoError(t, db.Write(dsA))
	dsB, err := db.Read(dsA.Name)
	assert.NoError(t, err)
	assert.True(t, mtx.Equal(dsA.Mtx, dsB.Mtx))
	assert.Equal(t, 8, len(dsB.Columns))
	assert.Equal(t, "TestReadWrite", dsB.Name)
}

func TestQuery(t *testing.T) {
	db, cleanup := NewTestDB(t)
	defer cleanup()
	mx1 := mtx.NewDense(2, 4, []float64{
		1.0, 1.0, 1.0, 1.0,
		2.0, 2.0, 2.0, 2.0,
	})
	mx2 := mtx.NewDense(3, 4, []float64{
		3.0, 3.0, 3.0, 3.0,
		2.0, 2.0, 2.0, 2.0,
		1.0, 1.0, 1.0, 1.0,
	})
	assert.NoError(t, db.Write(&types.Dataset{
		Name:    "mx1",
		Mtx:     mx1,
		Columns: []string{"A", "B", "C", "D"}}),
	)
	assert.NoError(t, db.Write(&types.Dataset{
		Name:    "mx2",
		Mtx:     mx2,
		Columns: []string{"E", "F", "G", "H"}}),
	)
	// Ensure multiple queries for the same dataset do not
	// return multiple rows
	q := types.Queries{
		&types.Query{Name: "mx1", Columns: []string{"A", "B", "B"}},
		&types.Query{Name: "mx1", Columns: []string{"B", "C", "D"}},
	}
	ds, err := db.Query(q)
	assert.NoError(t, err)
	mx := ds.Mtx
	r, c := mx.Dims()
	assert.Equal(t, 2, r)
	assert.Equal(t, 6, c)
	q = types.Queries{
		&types.Query{Name: "mx1", Columns: []string{"A", "B", "C"}},
		&types.Query{Name: "mx2", Columns: []string{"E", "F", "G"}},
	}
	ds, err = db.Query(q)
	mx = ds.Mtx
	assert.NoError(t, err)
	r, c = mx.Dims()
	assert.Equal(t, 5, r)
	assert.Equal(t, 6, c)
	assert.Equal(t, 6, len(ds.Columns))
	assert.Equal(t, 1.0, mx.At(0, ds.CPos("A")))
	assert.Equal(t, 2.0, mx.At(1, ds.CPos("A")))
	assert.Equal(t, 1.0, mx.At(0, ds.CPos("B")))
	assert.Equal(t, 2.0, mx.At(1, ds.CPos("B")))
	assert.Equal(t, 1.0, mx.At(0, ds.CPos("C")))
	assert.Equal(t, 2.0, mx.At(1, ds.CPos("C")))
	assert.Equal(t, 3.0, mx.At(0, ds.CPos("E")))
	assert.Equal(t, 2.0, mx.At(1, ds.CPos("E")))
	assert.Equal(t, 1.0, mx.At(2, ds.CPos("E")))
	assert.Equal(t, 3.0, mx.At(0, ds.CPos("F")))
	assert.Equal(t, 2.0, mx.At(1, ds.CPos("F")))
	assert.Equal(t, 1.0, mx.At(2, ds.CPos("F")))
	assert.Equal(t, 3.0, mx.At(0, ds.CPos("G")))
	assert.Equal(t, 2.0, mx.At(1, ds.CPos("G")))
	assert.Equal(t, 1.0, mx.At(2, ds.CPos("G")))
	_, err = db.Query(types.Queries{&types.Query{Name: "mx3"}})
	assert.Error(t, err, "not found")
	_, err = db.Query(types.Queries{&types.Query{Name: "mx1", Columns: []string{"H"}}})
	assert.Error(t, err, "not found")
}

func init() {
	rand.Seed(time.Now().Unix())
}
