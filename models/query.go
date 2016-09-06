package models

import (
	"net/url"
	"strconv"
	"time"
)

// Query is used to query the database
type Query struct {
	Start     time.Time         // Starting time of the entry
	End       time.Time         // Ending time of the entry
	Precision Precision         // Level of precision (value rollup) to apply
	Order     string            // Order of rows by time
	Match     map[string]string // Optional value to match on
}

// QueryFromURL builds a database query from a URL query string
// A full query with all options specified might look like:
// ?match=sport,Running&aggr=day&end=2016-Sep-03&start=2016-Aug-03
func QueryFromURL(u *url.URL) Query {
	values := u.Query()
	query := Query{}
	query.Start, _ = time.Parse("2006-Jan-02", values.Get("start"))
	query.End, _ = time.Parse("2006-Jan-02", values.Get("end"))
	if query.Start.IsZero() || query.End.IsZero() {
		query.End = time.Now()
		query.Start = query.End.AddDate(0, 0, -7) // Default to one week ago
	}
	val, _ := strconv.ParseInt(values.Get("precision"), 0, 64)
	query.Precision = Precision(val)
	if m := values.Get("match"); m != "" {
		if len(m) == 2 { // Only support a single criteria for the moment
			query.Match[string(m[0])] = string(m[1])
		}
	}
	query.Order = values.Get("order")
	return query
}
