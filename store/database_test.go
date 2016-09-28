package store

import (
	mtx "github.com/gonum/matrix/mat64"
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

func NewTestDB(t *testing.T) (*DB, func()) {
	f, err := ioutil.TempFile("/tmp", "fit-test")
	if err != nil {
		t.Error(err)
	}
	db, err := NewDB(f.Name())
	assert.NoError(t, err)
	return db, func() {
		db.Close()
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
	dsA := &Dataset{
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

func TestQueryFromArgs(t *testing.T) {
	args := []string{"D0,x,y", "D1,x"}
	assert.Equal(t, 2, len(QueryFromArgs(args)))
	assert.Equal(t, "D0", QueryFromArgs(args)[0].Name)
	assert.Equal(t, "x", QueryFromArgs(args)[0].Columns[0])
	assert.Equal(t, "y", QueryFromArgs(args)[0].Columns[1])
	assert.Equal(t, "D1", QueryFromArgs(args)[1].Name)
	assert.Equal(t, "x", QueryFromArgs(args)[1].Columns[0])
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
	assert.NoError(t, db.Write(&Dataset{
		Name:    "mx1",
		Mtx:     mx1,
		Columns: []string{"A", "B", "C", "D"}}),
	)
	assert.NoError(t, db.Write(&Dataset{
		Name:    "mx2",
		Mtx:     mx2,
		Columns: []string{"E", "F", "G", "H"}}),
	)
	q := []*Query{
		&Query{Name: "mx1", Columns: []string{"A", "B", "C"}},
		&Query{Name: "mx2", Columns: []string{"E", "F", "G"}},
	}
	ds, err := db.Query(q...)
	mx := ds.Mtx
	assert.NoError(t, err)
	r, c := mx.Dims()
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
	_, err = db.Query(&Query{Name: "mx3"})
	assert.Error(t, err, "not found")
	_, err = db.Query(&Query{Name: "mx1", Columns: []string{"H"}})
	assert.Error(t, err, "not found")
}

func init() {
	rand.Seed(time.Now().Unix())
}
