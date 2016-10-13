package types

import (
	mtx "github.com/gonum/matrix/mat64"
	"net/url"
	"strings"
)

// Query can be used to combine the results
// of multiple datasets into a single
// matrix of values. Queries can originate
// from the command line as arguments,
// a URL query string, or a JSON encoded
// payload.
//
// Command line arguments take the same
// form as URL encoding
//
// Text Specification:
// d=DS1,x,y&d=DS2,z,fuu&grouping=Duration,0,1m&fn=avg
//
type Query struct {
	Datasets []struct {
		Name    string   // Name of the dataset
		Columns []string // Columns within the dataset to query
	}
	Function *Function
	Grouping *Grouping
}

// Len returns the length of the Query
func (query Query) Len() int {
	return len(query.Datasets)
}

// Columns returns a flattened ordered
// array of Column names
func (query Query) Columns() []string {
	columns := make([]string, 0)
	for _, dataset := range query.Datasets {
		for _, column := range dataset.Columns {
			columns = append(columns, column)
		}
	}
	return columns
}

// String returns a valid URL query string
func (query Query) String() string {
	values := url.Values{}
	values.Add("fn", query.Function.Name)
	if query.Grouping != nil {
		values.Add("grouping", query.Grouping.String())
	}
	for _, dataset := range query.Datasets {
		args := make([]string, len(dataset.Columns)+1)
		args[0] = dataset.Name
		for i, column := range dataset.Columns {
			args[i+1] = column
		}
		values.Add("q", strings.TrimRight(strings.Join(args, ","), ","))
	}
	return values.Encode()
}

// Apply returns a new modified matrix based on the query
func (query Query) Apply(mx *mtx.Dense) *mtx.Dense {
	if query.Grouping != nil {
		return query.Function.Apply(query.Grouping.Group(mx))
	}
	return mx
}

// NewQuery constructs a Query from the provided
// args and optionally specified function.
// If function is specified the query returns
// aggregated
func NewQuery(args []string, function, grouping string) *Query {
	query := &Query{
		Datasets: make([]struct {
			Name    string
			Columns []string
		}, len(args)),
		Function: &Function{
			Name: function,
		},
	}
	if grouping != "" {
		query.Grouping = NewGrouping(grouping)
	}
	for i, arg := range args {
		split := strings.Split(arg, ",")
		if len(split) >= 1 {
			query.Datasets[i].Name = split[0]
		}
		if len(split) > 1 {
			query.Datasets[i].Columns = split[1:]
		}
	}
	return query
}

// NewQueryQS constructs a query from a url.URL
func NewQueryQS(u *url.URL) *Query {
	var args []string
	query := u.Query()
	if q, ok := query["q"]; ok {
		args = q
	}
	return NewQuery(args, query.Get("fn"), query.Get("grouping"))
}
