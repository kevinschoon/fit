package server

import (
	"net/url"
	"strconv"
	//"strings"
	"github.com/kevinschoon/gofit/models"
	"time"
)

// Query is used to query the database
type Query struct {
	Name      string            // Name of the series
	Start     time.Time         // Starting time of the entry
	End       time.Time         // Ending time of the entry
	Precision models.Precision  // Level of precision (value rollup) to apply
	Match     map[string]string // Optional value to match on
	X         models.Key        // X axis
	Y         models.Key        // Y axis
}

// QueryFromURL builds a database query from a URL query string
// A full query with all options specified might look like:
// ?match=sport,Running&aggr=day&end=2016-Sep-03&start=2016-Aug-03
func QueryFromURL(u *url.URL) *Query {
	values := u.Query()
	query := &Query{
		Match: make(map[string]string),
	}
	query.Start, _ = time.Parse("2006-Jan-02", values.Get("start"))
	query.End, _ = time.Parse("2006-Jan-02", values.Get("end"))
	if query.End.IsZero() { // Default to showing all data until the current time.
		query.End = time.Now()
	}
	val, _ := strconv.ParseInt(values.Get("precision"), 0, 64)
	query.Precision = models.Precision(val)
	X, _ := strconv.ParseInt(values.Get("X"), 0, 64)
	Y, _ := strconv.ParseInt(values.Get("Y"), 0, 64)
	if X == 0 && Y == 0 {
		Y = 1
	}
	query.X = models.Key(X)
	query.Y = models.Key(Y)
	return query
}
