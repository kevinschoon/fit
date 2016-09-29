package csv

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

var Simple string = `
"","time","LakeHuron"
"1",1875,580.38
"2",1876,581.86
"3",1877,580.97
"4",1878,580.8
"5",1879,579.79
`

func TestCSVLoader(t *testing.T) {
	c, err := New(strings.NewReader(Simple), nil)
	assert.NoError(t, err)
	values, err := c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(values))
	assert.Equal(t, 580.38, values[2])
	values, err = c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 581.86, values[2])
	values, err = c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 580.97, values[2])
	values, err = c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 580.8, values[2])
	values, err = c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 579.79, values[2])
	_, err = c.Next()
	assert.Error(t, err)
	assert.Equal(t, err, io.EOF)
}
