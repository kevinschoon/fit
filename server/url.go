package server

import (
	"fmt"
	"net/url"
	"time"
)

// URLBuilder modifies the query string of an existing URL
type URLBuilder struct {
	URL        *url.URL // URL of the request
	Collection string   // Name of the collection
}

// Copy copies the existing URL
func (u URLBuilder) Copy() *url.URL {
	return &url.URL{
		Scheme: u.URL.Scheme,
		Host:   u.URL.Host,
		Path:   u.URL.Path,
	}
}

// Chart returns a URL to generate the corresponding chart
//for a given query
func (u URLBuilder) Chart() string {
	values := u.URL.Query()
	copied := u.Copy()
	copied.Path = fmt.Sprintf("%s/chart", u.Collection)
	fmt.Println(copied.Path)
	copied.RawQuery = values.Encode()
	return copied.String()
}

// ByRange returns a URL with a time spanning a week, month, or year
func (u URLBuilder) ByRange(key string) string {
	values := u.URL.Query()
	copied := u.Copy()
	now := time.Now()
	values.Set("end", now.Format("2006-Jan-02"))
	switch key {
	case "week":
		values.Set("start", now.AddDate(0, 0, -7).Format("2006-Jan-02"))
	case "month":
		values.Set("start", now.AddDate(0, -1, 0).Format("2006-Jan-02"))
	case "year":
		values.Set("start", now.AddDate(-1, 0, 0).Format("2006-Jan-02"))
	}
	copied.RawQuery = values.Encode()
	return copied.String()
}

// ByTime returns a URL with a query based on an existing time (like TCX StartTime)
func (u URLBuilder) ByTime(t time.Time) string {
	values := u.URL.Query()
	copied := u.Copy()
	values.Set("end", t.Format("2006-Jan-02"))
	values.Set("start", t.Format("2006-Jan-02"))
	values.Set("precision", "3")
	copied.RawQuery = values.Encode()
	return copied.String()
}

// ByAggr returns a URL with the specified precision level
func (u URLBuilder) ByPrecision(key string) string {
	values := u.URL.Query()
	copied := u.Copy()
	values.Set("precision", key)
	copied.RawQuery = values.Encode()
	return copied.String()
}

// ByMatch returns a URL specifying a match criteria
func (u URLBuilder) ByMatch(key, value string) string {
	values := u.URL.Query()
	copied := u.Copy()
	if value == "all" {
		values.Del("match")
	}
	values.Set("match", fmt.Sprintf("%s,%s", key, value))
	copied.RawQuery = values.Encode()
	return copied.String()
}
