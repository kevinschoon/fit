package models

import "time"

type Precision int
type Key int

const (
	Days Precision = iota
	Months
	Years
)

type Value struct {
	Name  string
	Value float64
}

// Series is a group of values aggregated by time
type Series struct {
	index  int
	Time   time.Time
	Values [][]Value
}

func (series Series) Len() int {
	return len(series.Values)
}

func (series Series) Less(i, j int) bool {
	return series.Values[i][series.index].Value < series.Values[j][series.index].Value
}

func (series Series) Swap(i, j int) {
	series.Values[i], series.Values[j] = series.Values[j], series.Values[i]
}

func (series *Series) SortBy(k Key) {
	series.index = int(k)
}

func (series Series) Get(i int, k Key) Value {
	return series.Values[i][int(k)]
}
func (series Series) GetAll(i int) []Value {
	return series.Values[i]
}

func (series *Series) Add(values []Value) {
	series.Values = append(series.Values, values)
}

type Collection struct {
	Series []*Series
	Name   string
}

func (c Collection) Len() int {
	return len(c.Series)
}

// Add enters a new series of values into the collection
// By default values are aggregated and stored per second
func (c *Collection) Add(start time.Time, values []Value) {
	if len(c.Series) == 0 {
		c.Series = []*Series{
			&Series{
				Time:   start,
				Values: [][]Value{values},
			},
		}
		return
	}
	if series := c.Find(start); series != nil {
		series.Add(values)
		//series.Values = append(series.Values, values)
	} else {
		c.Series = append(c.Series, &Series{
			Time:   start,
			Values: [][]Value{values},
		})
	}
}

// Find searches for a series matching the provided start time
func (c *Collection) Find(start time.Time) *Series {
	for _, series := range c.Series {
		if series.Time.Unix() == start.Unix() {
			return series
		}
	}
	return nil
}

// Names returns value names in the collection of Series
func (c *Collection) Names() (names []string) {
	if c.Len() > 0 {
		for _, value := range c.Series[0].GetAll(0) {
			names = append(names, value.Name)
		}
	}
	return names
}

// Rollup aggregates series by the specified precision
func (c *Collection) RollUp(precision Precision) {
	collection := Collection{
		Name: c.Name,
	}
	aggr := make(map[int][]*Series)
	for _, series := range c.Series {
		var key int
		switch precision {
		case Years:
			key = series.Time.Year()
		case Months:
			key = (series.Time.Year() * 12) + int(series.Time.Month())
		case Days:
			key = (series.Time.Year() * 12) + (int(series.Time.Month()) * 31) + series.Time.Day()
		}
		if _, ok := aggr[key]; !ok {
			aggr[key] = make([]*Series, 0)
		}
		aggr[key] = append(aggr[key], series)
	}
	for _, series := range aggr {
		first := series[0]
		if len(series) > 1 {
			for _, series := range series[1:] {
				for _, values := range series.Values {
					first.Add(values)
				}
			}
		}
		collection.Series = append(collection.Series, first)
	}
	*c = collection
}

/*
func (rows Rows) RollUp(precision Precision) Rows {
	if precision == None {
		return rows
	}
	buckets := make(map[int]Rows)
	for _, row := range rows {
		var key int
		switch precision {
		case Years:
			key = row.Time.Year()
		case Months:
			key = (row.Time.Year() * 12) + int(row.Time.Month())
		case Days:
			key = (row.Time.Year() * 12) + (int(row.Time.Month()) * 31) + row.Time.Day()
		}
		if _, ok := buckets[key]; !ok {
			buckets[key] = make(Rows, 0)
		}
		buckets[key] = append(buckets[key], row)
	}
	newRows := Rows{}
	for _, bucket := range buckets {
		first := bucket[0]
		if len(bucket) > 1 {
			bucket = bucket[1:]
			for _, row := range bucket {
				for i, value := range row.Values {
					first.Values[i] += value
				}
			}
		}
		newRows = append(newRows, first)
	}
	return newRows
}
*/
