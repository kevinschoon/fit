package store

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestQueries(t *testing.T) {
	args := []string{"D0,x,y", "D1,x,z"}
	queries := NewQueries(args)
	assert.Equal(t, 2, len(queries))
	assert.Equal(t, "D0", queries[0].Name)
	assert.Equal(t, "D1", queries[1].Name)
	columns := queries.Columns()
	assert.Equal(t, 4, len(columns))
	assert.Equal(t, "x", columns[0])
	assert.Equal(t, "y", columns[1])
	assert.Equal(t, "x", columns[2])
	assert.Equal(t, "z", columns[3])
	assert.Equal(t, "q=D0%2Cx%2Cy&q=D1%2Cx%2Cz", queries.QueryStr())
}

func TestQueriesFromQS(t *testing.T) {
	u, err := url.Parse("http://localhost?q=Fuu,x&q=Bar,y,z")
	assert.NoError(t, err)
	queries := NewQueriesFromQS(u)
	assert.Equal(t, 2, queries.Len())
	columns := queries.Columns()
	assert.Len(t, columns, 3)
	assert.Equal(t, "x", columns[0])
	assert.Equal(t, "y", columns[1])
	assert.Equal(t, "z", columns[2])

}
