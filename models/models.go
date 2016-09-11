package models

import "time"

type Precision int
type Key int

const (
	Days Precision = iota
	Months
	Years
	None
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

func (series Series) Sum(k Key) (sum float64) {
	for _, values := range series.Values {
		if len(values) > 0 {
			if len(values) >= int(k) {
				sum += values[int(k)].Value
			}
		}
	}
	return sum
}

func (series Series) SumAll() (sums []float64) {
	if len(series.Values) > 0 {
		sums = make([]float64, len(series.Values[0]))
		for _, values := range series.Values {
			for i, value := range values {
				if len(sums) >= i {
					sums[i] += value.Value
				}
			}
		}
	}
	return sums
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

// Name returns the name of a given Series key
func (c *Collection) GetName(k Key) (name string) {
	names := c.Names()
	if int(k) <= len(names) {
		name = names[int(k)]
	}
	return name
}

// Flat flattens each Series of values into their own Series
func (c *Collection) Flat() {
	collection := Collection{
		Name: c.Name,
	}
	for _, series := range c.Series {
		for _, values := range series.Values {
			collection.Series = append(collection.Series, &Series{
				Values: [][]Value{
					values,
				},
			})
		}
	}
	*c = collection
}

// Rollup aggregates series by the specified precision
func (c *Collection) RollUp(precision Precision) {
	if precision == None {
		c.Flat()
		return
	}
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
