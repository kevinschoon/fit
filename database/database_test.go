package database

import (
	"github.com/kevinschoon/gofit/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

var cleanup bool = false

func Series(count int) []*models.Series {
	series := make([]*models.Series, 1)
	series[0] = models.NewSeries([]string{"V1", "V2"})
	series[0].Name = "TestSeries"
	start := time.Date(2016, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < count; i++ {
		series[0].Add(start, []models.Value{
			models.Value(rand.Float64()),
			models.Value(rand.Float64()),
		})
		start = start.Add(1 * time.Second)
	}
	return series
}

func newDB(t *testing.T) (*DB, func()) {
	f, err := ioutil.TempFile("/tmp", "gofit")
	if err != nil {
		t.Error(err)
	}
	db, err := New(f.Name(), true)
	assert.NoError(t, err)
	fn := func() {
		db.Close()
		if cleanup {
			os.Remove(f.Name())
		}
	}
	return db, fn
}

func TestDBSeriesReadWrite(t *testing.T) {
	db, cleanup := newDB(t)
	defer cleanup()
	series := Series(10000)
	assert.Equal(t, 1, len(series))
	series[0].Name = "SeriesReadWrite"
	assert.Equal(t, 3, len(series[0].Keys))
	assert.Equal(t, models.Key(0), series[0].Keys["time"])
	assert.Equal(t, models.Key(1), series[0].Keys["V1"])
	assert.Equal(t, models.Key(2), series[0].Keys["V2"])
	assert.NoError(t, db.WriteSeries(series))
	series, err := db.Series()
	assert.Equal(t, 1, len(series))
	series, err = db.ReadSeries("SeriesReadWrite", time.Time{}, time.Now())
	assert.NoError(t, err)
	// Series should be loaded at original 1/s interval
	assert.Equal(t, 10000, len(series))
	assert.Equal(t, 3, len(series[0].Keys))
	assert.Equal(t, models.Key(0), series[0].Keys["time"])
	assert.Equal(t, models.Key(1), series[0].Keys["V1"])
	assert.Equal(t, models.Key(2), series[0].Keys["V2"])
}

func init() {
	rand.Seed(time.Now().Unix())
}
