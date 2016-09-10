package database

import (
	t "github.com/kevinschoon/gofit/models/tcx"
	"github.com/kevinschoon/tcx"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

var cleanup bool = false

func getTCX(count int) *t.TCX {
	t := &t.TCX{
		Acts: make(tcx.Acts, count),
	}
	for x := 0; x < count; x++ {
		t.Acts[x] = tcx.Activity{
			StartTime: time.Now().AddDate(0, 0, -x),
			Laps:      make([]tcx.Lap, 5),
		}
		for i := 0; i < 5; i++ {
			t.Acts[x].Laps[i] = tcx.Lap{
				TotalTime: rand.Float64(),
				Dist:      rand.Float64(),
			}
		}
	}
	return t
}

func newDB(t *testing.T) (*DB, func()) {
	f, err := ioutil.TempFile("/tmp", "gofit")
	if err != nil {
		t.Error(err)
	}
	db, err := New(f.Name())
	assert.NoError(t, err)
	fn := func() {
		db.Close()
		if cleanup {
			os.Remove(f.Name())
		}
	}
	return db, fn
}

func TestDatabase(t *testing.T) {
	db, cleanup := newDB(t)
	defer cleanup()
	collection := getTCX(1000).Load()
	assert.NoError(t, db.Write("stuff", collection))
	collection, err := db.Read("stuff", time.Time{}, time.Now())
	assert.NoError(t, err)
	assert.Equal(t, collection.Len(), 1000)
}

func init() {
	rand.Seed(time.Now().Unix())
}
