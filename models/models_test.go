package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestSeriesSort(t *testing.T) {
	series := &Series{}
	series.Add([]Value{
		Value{
			Name:  "V1-0",
			Value: 1.0,
		},
		Value{
			Name:  "V1-1",
			Value: 1.0,
		},
	})
	series.Add([]Value{
		Value{
			Name:  "V2-0",
			Value: 1.0,
		},
		Value{
			Name:  "V2-1",
			Value: 2.0,
		},
	})
	assert.Equal(t, series.Len(), 2)
	series.SortBy(Key(1))
	sort.Sort(sort.Reverse(series))
	assert.Equal(t, "V2-0", series.GetAll(0)[0].Name)
	assert.Equal(t, "V2-1", series.GetAll(0)[1].Name)
	assert.Equal(t, "V1-0", series.GetAll(1)[0].Name)
	assert.Equal(t, "V1-1", series.GetAll(1)[1].Name)
}

func TestCollectionAdd(t *testing.T) {
	collection := Collection{}
	start := time.Now()
	for time.Now().Sub(start) < 2*time.Second {
		collection.Add(time.Now(), []Value{
			Value{
				Value: rand.Float64(),
			},
			Value{
				Value: rand.Float64(),
			},
		})
	}
	// This could occasionally fail if the test begins
	// at the beginning of a second in which case the
	// length should be 2.
	assert.Equal(t, 3, len(collection.Series))
	fmt.Println(len(collection.Series))
	// Add a new set of values during the timespan of the previous operation
	// the method should find the previous series and append it there
	collection.Add(start.Add(1*time.Second), []Value{
		Value{
			Value: rand.Float64(),
		},
	})
	assert.Equal(t, 3, len(collection.Series))
	// Add another entry several seconds after all the previous values
	// ensuring a new Series is created
	collection.Add(time.Now().Add(10*time.Second), []Value{
		Value{
			Value: rand.Float64(),
		},
	})
	assert.Equal(t, 4, len(collection.Series))
}

func TestCollectionRollup(t *testing.T) {
	collection := &Collection{}
	previous := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 23; i++ {
		// Increase the month each iteration
		previous = previous.AddDate(0, 1, 0)
		collection.Add(previous, []Value{
			Value{
				Value: rand.Float64(),
			},
			Value{
				Value: rand.Float64(),
			},
		})
	}
	assert.Equal(t, 23, collection.Len())
	// Rollup the collection by years, starting from January 1, 23 months
	// should be aggregated into two year long series
	collection.RollUp(Years)
	assert.Equal(t, 2, collection.Len())
}

func init() {
	rand.Seed(time.Now().Unix())
}
