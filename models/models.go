package models

import (
	"strconv"
	"time"
)

// Key is an identifier to represent a value within a series
// Keys can be used to associate a value with a specific
// type and to represent the index of a set of values
type Key int

// Keys hold the name and index for Value
type Keys map[string]Key

func (k Keys) Names() []string {
	names := make([]string, len(k))
	for name, key := range k {
		names[int(key)] = name
	}
	return names
}

// Value represents a single datapoint
type Value float64

func (value Value) Float64() float64 {
	return float64(value)
}

func (value Value) String() string {
	return strconv.FormatFloat(float64(value), 'E', -1, 64)
}

func (value Value) Time() time.Time {
	return time.Unix(int64(value), 64).UTC()
}

// Sortable implements the sort.Sort interface for an array of Value
type Sortable []Value

func (s Sortable) Len() int           { return len(s) }
func (s Sortable) Less(i, j int) bool { return s[i] < s[j] }
func (s Sortable) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Values represents a collection of Value arrays
type Values [][]Value

// Series is an append-only collection of Values grouped by time
// A series can apply different functions to the data combined within
type Series struct {
	values Values `json:"-"` // Arbitrarily sized collection of values
	index  int
	Name   string
	Keys   Keys
}

// Implements interface for sort.Sort
func (s Series) Len() int           { return len(s.values) }
func (s Series) Less(i, j int) bool { return s.values[i][0] < s.values[j][0] }
func (s Series) Swap(i, j int)      { s.values[i], s.values[j] = s.values[j], s.values[i] }

// Exists checks to see if values exist at a given index
// optionally it checks if a key exists at said index
func (series Series) Exists(i int, k Key, check bool) bool {
	if len(series.values) >= i+1 { // Series exists at index i
		if !check { // Not checking for key
			return true
		}
	} else { // Series out of range
		return false
	}
	if len(series.values[i]) >= int(k)+1 { // Key exists at series index i
		return true
	}
	return false // Key does not exist at series index i
}

// Value returns the value for k at the given index
// if the value exists, otherwise it returns 0
// TODO: Rework to accept a key rather than string
func (series Series) Value(i int, k string) Value {
	if key, ok := series.Keys[k]; ok {
		if series.Exists(i, key, true) {
			return series.values[i][int(key)]
		}
	}
	return Value(0) // Return zero values for missing data
}

// Values dumps all of the values for the given key
func (series Series) Values(key Key) []Value {
	values := make([]Value, len(series.values))
	for i := 0; i < len(series.values); i++ {
		if series.Exists(i, key, true) {
			values[i] = series.values[i][int(key)]
		} else {
			values[i] = Value(0)
		}
	}
	return values
}

// Start returns the time of the first value within the series
func (series Series) Start() time.Time {
	return series.Value(0, "time").Time()
}

// End returns the time of the last value within the series
func (series Series) End() time.Time {
	if len(series.values) >= 1 {
		return series.Value(len(series.values)-1, "time").Time()
	}
	return series.Start()
}

// Add enters a new array of Value into the series
func (series *Series) Add(start time.Time, values []Value) {
	start = start.UTC()
	if start.IsZero() {
		start = time.Now()
	}
	v := []Value{Value(start.Unix())}
	for _, value := range values {
		v = append(v, value)
	}
	series.values = append(series.values, v)
}

// Next is used to iterate over a series
func (series *Series) Next() (values []Value) {
	if series.Exists(series.index, Key(0), false) {
		values = series.values[series.index]
		series.index++ // TODO: Mutex
	} else {
		series.index = 0
	}
	return values
}

// Dump returns all of the values in the Series
func (series Series) Dump() Values {
	return series.values
}

// Import replaces all current values with those provided
func (series *Series) Import(values Values) {
	series.values = values
}

// New series creates a new Series
func NewSeries(name string, columns []string) *Series {
	series := &Series{
		Name: name,
		Keys: map[string]Key{
			"time": Key(0),
		},
	}
	for i := 1; i < len(columns)+1; i++ {
		if columns[i-1] == "time" {
			series.Keys["_time"] = Key(i)
			continue
		}
		series.Keys[columns[i-1]] = Key(i)
	}
	return series
}

// Resize takes an array of series and arranges
// the underlying values based on the specified duration
// TODO: Ensure that all series are of the same name
// and share the same keys
func Resize(input []*Series, aggr time.Duration) (output []*Series) {
	var current *Series
	for _, series := range input {
		// We are at the first Series
		if len(output) == 0 {
			// Shallow copy Series
			current = Copy(series)
			// Set the values of the current series to
			// the next set of Values in the array
			current.values = Values{series.Next()}
			// Append the current Series to the output
			output = append(output, current)
		}
		// Iterate the Series until the end of it's values
		for values := series.Next(); values != nil; values = series.Next() {
			// Get the duration of all the values within the Series
			duration := current.End().Sub(current.Start()) + values[0].Time().Sub(current.End())
			// The duration of the current series exceeds the
			// the aggregation threshold or is zero.
			if duration >= aggr || aggr == time.Duration(0) {
				// Shallow copy the series and set it current
				current = Copy(series)
				// Append the new Series to the output
				output = append(output, current)
			}
			// Add this series of Values to the current Series
			current.values = append(current.values, values)
		}
	}
	return output
}

// Copy performs a shallow copy of the given series
func Copy(in *Series) *Series {
	series := &Series{
		Name: in.Name,
		Keys: make(map[string]Key),
	}
	for name, key := range in.Keys {
		series.Keys[name] = Key(int(key))
	}
	return series
}

// Select returns all of the values for a given
// key across an array of Series
func Select(key Key, series []*Series) (values []Value) {
	for _, s := range series {
		for _, value := range s.Values(key) {
			values = append(values, value)
		}
	}
	return values
}
