package models

import "time"

type Precision int

const (
	Days Precision = iota
	Months
	Years
	None
)

// Datapoint is a single pair of XY values
type Datapoint struct {
	X float64
	Y float64
}

// Datapoints is an array of XY values
type Datapoints []Datapoint

func (dps Datapoints) Len() int {
	return len(dps)
}

func (dps Datapoints) Value(i int) float64 {
	return dps[i].Y
}

func (dps Datapoints) XY(i int) (x, y float64) {
	return dps[i].X, dps[i].Y
}

type Value float64

type Row struct {
	Time   time.Time
	Values []Value
}

type Rows []Row

func (rows Rows) Len() int {
	return len(rows)
}

func (rows Rows) Less(i, j int) bool {
	return rows[i].Time.Unix() < rows[j].Time.Unix()
}

func (rows Rows) Swap(i, j int) {
	rows[i], rows[j] = rows[j], rows[i]
}

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

type Series interface {
	Columns() []string
	Rows() Rows
	Pts(string) Datapoints // TODO: Change to concrete types (int)
}
