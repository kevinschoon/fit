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

type Value float64

type XYs []struct{ X, Y float64 }

func (xys XYs) Len() int {
	return len(xys)
}

func (xys XYs) XY(i int) (float64, float64) {
	return xys[i].X, xys[i].Y
}

func (xys XYs) Less(i, j int) bool {
	return xys[i].X < xys[j].Y
}

func (xys XYs) Swap(i, j int) {
	xys[i], xys[j] = xys[j], xys[i]
}

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

type Series struct {
	Columns []string
	Rows    Rows
}

func (series Series) Pts(y Key) XYs {
	xys := make(XYs, len(series.Rows))
	for i := range xys {
		xys[i].X = float64(series.Rows[i].Time.Unix())
		if len(series.Rows[i].Values) < int(y) {
			panic("Invalid series index")
		}
		xys[i].Y = float64(series.Rows[i].Values[int(y)])
	}
	return xys
}

type Serieser interface {
	Series() *Series
	Types() []interface{}
}
