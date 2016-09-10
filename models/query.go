package models

import (
	"net/url"
	"strconv"
	//"strings"
	"time"
)

// Query is used to query the database
type Query struct {
	Name      string            // Name of the series
	Start     time.Time         // Starting time of the entry
	End       time.Time         // Ending time of the entry
	Precision Precision         // Level of precision (value rollup) to apply
	Match     map[string]string // Optional value to match on
}

// QueryFromURL builds a database query from a URL query string
// A full query with all options specified might look like:
// ?match=sport,Running&aggr=day&end=2016-Sep-03&start=2016-Aug-03
func QueryFromURL(u *url.URL) Query {
	values := u.Query()
	query := Query{
		Match: make(map[string]string),
	}
	query.Start, _ = time.Parse("2006-Jan-02", values.Get("start"))
	query.End, _ = time.Parse("2006-Jan-02", values.Get("end"))
	if query.End.IsZero() { // Default to showing all data until the current time.
		query.End = time.Now()
	}
	val, _ := strconv.ParseInt(values.Get("precision"), 0, 64)
	query.Precision = Precision(val)

	return query
}

/*
	if query.Start.IsZero() || query.End.IsZero() {
		query.End = time.Now()
		query.Start = query.End.AddDate(0, 0, -7) // Default to one week ago
	}
*/
/*
	if m := values.Get("match"); m != "" {
		split := strings.Split(m, ",")
		if len(split) == 2 { // TODO: Support multiple criteria
			query.Match[split[0]] = split[1]
		}
	}

	xval, _ := strconv.ParseInt(values.Get("X"), 0, 64)
	yval, _ := strconv.ParseInt(values.Get("Y"), 0, 64)
	query.X = Key(xval)
	if yval == 0 {
		yval++
	}
	query.Y = Key(yval)
*/
