package loader

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
	c, err := NewCSV(strings.NewReader(Simple))
	assert.NoError(t, err)
	assert.Len(t, c.Columns, 3)
	rows, cols := c.Dims()
	assert.Equal(t, 5, rows)
	assert.Equal(t, 3, cols)
	values, err := c.Row()
	assert.NoError(t, err)
	assert.Len(t, values, 3)
	assert.Equal(t, "1", values[0])
	assert.Equal(t, "1875", values[1])
	assert.Equal(t, "580.38", values[2])
	values, err = c.Row()
	assert.NoError(t, err)
	assert.Len(t, values, 3)
	assert.Equal(t, "2", values[0])
	assert.Equal(t, "1876", values[1])
	assert.Equal(t, "581.86", values[2])
	values, err = c.Row()
	assert.NoError(t, err)
	assert.Equal(t, "3", values[0])
	assert.Equal(t, "1877", values[1])
	assert.Equal(t, "580.97", values[2])
	values, err = c.Row()
	assert.NoError(t, err)
	assert.Equal(t, "4", values[0])
	assert.Equal(t, "1878", values[1])
	assert.Equal(t, "580.8", values[2])
	values, err = c.Row()
	assert.NoError(t, err)
	assert.Equal(t, "5", values[0])
	assert.Equal(t, "1879", values[1])
	assert.Equal(t, "579.79", values[2])
	_, err = c.Row()
	assert.Error(t, err)
	assert.Equal(t, err, io.EOF)
}
