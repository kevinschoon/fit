package models

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestSeriesSort(t *testing.T) {
	series := NewSeries([]string{"V1", "V2"})
	series.Add(time.Now().AddDate(-1, 0, 0), []Value{Value(2.0), Value(1.0)})
	series.Add(time.Now(), []Value{Value(1.0), Value(2.0)})
	assert.Equal(t, Value(2.0), series.Value(0, "V1"))
	assert.Equal(t, Value(1.0), series.Value(0, "V2"))
	assert.Equal(t, Value(1.0), series.Value(1, "V1"))
	assert.Equal(t, Value(2.0), series.Value(1, "V2"))
	sort.Sort(sort.Reverse(series))
	assert.Equal(t, Value(1.0), series.Value(0, "V1"))
	assert.Equal(t, Value(2.0), series.Value(0, "V2"))
	assert.Equal(t, Value(2.0), series.Value(1, "V1"))
	assert.Equal(t, Value(1.0), series.Value(1, "V2"))
}

func TestSeriesExists(t *testing.T) {
	series := NewSeries([]string{"V1"})
	series.Add(time.Now(), []Value{Value(1.0)})
	series.Add(time.Now(), []Value{Value(2.0)})
	series.Add(time.Now(), []Value{Value(1.0)})
	assert.True(t, series.Exists(0, Key(0), true))
	assert.True(t, series.Exists(0, Key(1), true))
	assert.True(t, series.Exists(1, Key(0), true))
	assert.False(t, series.Exists(0, Key(2), true))
	assert.False(t, series.Exists(3, Key(0), false))
	assert.False(t, series.Exists(3, Key(1), true))
}

func TestNext(t *testing.T) {
	series := NewSeries([]string{"V1"})
	for i := 0; i < 10; i++ {
		series.Add(time.Now(), []Value{Value(rand.Float64())})
	}
	for i := 0; i < 10; i++ {
		assert.Equal(t, i, series.index)
		series.Next()
	}
	assert.Nil(t, series.Next())
	assert.Equal(t, 0, series.index)
}

func TestResize(t *testing.T) {
	series := make([]*Series, 1)
	start := time.Date(2016, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	series[0] = NewSeries([]string{"V1", "V2"})
	for i := 0; i < 60; i++ {
		series[0].Add(start, []Value{
			Value(rand.Float64()),
			Value(rand.Float64()),
		})
		start = start.Add(1 * time.Second)
	}
	assert.Equal(t, 60, len(Resize(series, 1*time.Second)))
	assert.Equal(t, 30, len(Resize(series, 2*time.Second)))
	assert.Equal(t, 1, len(Resize(series, 60*time.Second)))
	assert.Equal(t, 1, len(Resize(series, 1*time.Hour)))

	// Wide series over 24 hours
	series = make([]*Series, 24)
	start = time.Date(2016, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	for s := 0; s < 24; s++ {
		series[s] = NewSeries([]string{"V1", "V2"})
		for i := 0; i < 3600; i++ {
			series[s].Add(start, []Value{
				Value(rand.Float64()),
				Value(rand.Float64()),
			})
			start = start.Add(1 * time.Second)
		}
	}
	assert.Equal(t, 24, len(series))
	assert.Equal(t, 3600, len(series[0].values))
	assert.Equal(t, 86400, len(Resize(series, 1*time.Second)))
	assert.Equal(t, 24, len(Resize(series, 1*time.Hour)))
	assert.Equal(t, 1440, len(Resize(series, 1*time.Minute)))

	// Single large series
	series = make([]*Series, 1)
	start = time.Date(2016, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	series[0] = NewSeries([]string{"V1", "V2"})
	for i := 0; i < 86400; i++ {
		series[0].Add(start, []Value{
			Value(rand.Float64()),
			Value(rand.Float64()),
		})
		start = start.Add(1 * time.Second)
	}
	assert.Equal(t, 24, len(Resize(series, 1*time.Hour)))

	// Random unsorted times
	series = make([]*Series, 1)
	start = time.Date(2016, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	series[0] = NewSeries([]string{"V1", "V2"})
	for i := 0; i < 30; i++ {
		series[0].Add(start, []Value{
			Value(rand.Float64()),
			Value(rand.Float64()),
		})
		start = start.Add(time.Duration(rand.Intn(5)) * time.Minute)
	}
	assert.True(t, len(Resize(series, 1*time.Minute)) < 30)
	assert.True(t, len(Resize(series, 1*time.Minute)) > 1)
}

func TestSeriesFunction(t *testing.T) {
	series := make([]*Series, 10)
	start := time.Date(2016, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	// Add 1000 1 minute intervals
	for s := 0; s < 10; s++ {
		series[s] = NewSeries([]string{"V1", "V2"})
		for i := 0; i < 100; i++ {
			series[s].Add(start, []Value{
				Value(1.0),
				Value(2.0),
			})
			start = start.Add(1 * time.Minute)
		}
	}
	// Aggregate by day
	assert.Equal(t, 1, len(Resize(series, 24*time.Hour)))
	assert.Equal(t, 100, len(series[0].values))
	assert.Equal(t, 1, len(Sum(series[0]).values))
	assert.Equal(t, Value(100), Sum(series[0]).Value(0, "V1"))
	assert.Equal(t, Value(200), Sum(series[0]).Value(0, "V2"))
	assert.Equal(t, series[0].Start(), Sum(series[0]).Start())
	assert.Equal(t, len(series[0].values[0]), len(Sum(series[0]).values[0]))
	assert.Equal(t, 10, len(series))
	assert.Equal(t, 10, len(Apply(series, Sum)))
	for _, series := range Apply(series, Sum) {
		assert.Equal(t, Value(100), series.Value(0, "V1"))
		assert.Equal(t, Value(200), series.Value(0, "V2"))
	}
	series = make([]*Series, 10)
	start = time.Date(2016, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	series = []*Series{NewSeries([]string{"V1"})}
	series[0].Add(start, []Value{
		Value(1.2),
	})
	series[0].Add(start, []Value{
		Value(1.0),
	})
	series[0].Add(start, []Value{
		Value(2.0),
	})
	series[0].Add(start, []Value{
		Value(1.8),
	})
	series[0].Add(start, []Value{
		Value(1.6),
	})
	series[0].Add(start, []Value{
		Value(1.4),
	})
	assert.Equal(t, Value(1.5), Avg(series[0]).Value(0, "V1"))
	assert.Equal(t, Value(1.0), Min(series[0]).Value(0, "V1"))
	assert.Equal(t, Value(2.0), Max(series[0]).Value(0, "V1"))
}

func init() {
	rand.Seed(time.Now().Unix())
}
