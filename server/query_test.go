package server

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestQueries(t *testing.T) {
	u, err := url.Parse("http://localhost?x=Fuu,x&y=Bar,y,z")
	assert.NoError(t, err)
	queries := XYQueries(u)
	assert.Equal(t, 2, queries.Len())
	columns := queries.Columns()
	assert.Len(t, columns, 3)
	assert.Equal(t, "x", columns[0])
	assert.Equal(t, "y", columns[1])
	assert.Equal(t, "z", columns[2])
}
