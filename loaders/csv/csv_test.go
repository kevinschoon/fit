package csv

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
	"time"
)

var Simple string = `
"","time","LakeHuron"
"1",1875,580.38
"2",1876,581.86
"3",1877,580.97
"4",1878,580.8
"5",1879,579.79
`

var CustomDate string = `
dt,LandAverageTemperature,LandAverageTemperatureUncertainty,LandMaxTemperature,LandMaxTemperatureUncertainty,LandMinTemperature,LandMinTemperatureUncertainty,LandAndOceanAverageTemperature,LandAndOceanAverageTemperatureUncertainty
1750-01-01,3.0340000000000003,3.574,,,,,,
1750-02-01,3.083,3.702,,,,,,
1750-03-01,5.626,3.076,,,,,,
1750-04-01,8.49,2.451,,,,,,
1750-05-01,11.573,2.072,,,,,,
`

func TestCSV(t *testing.T) {
	c, err := New("", &Options{
		Reader: strings.NewReader(Simple),
	})
	assert.NoError(t, err)
	values, err := c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(values))
	assert.Equal(t, 580.38, values[2].Float64())
	values, err = c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 581.86, values[2].Float64())
	values, err = c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 580.97, values[2].Float64())
	values, err = c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 580.8, values[2].Float64())
	values, err = c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 579.79, values[2].Float64())
	_, err = c.Next()
	assert.Error(t, err)
	assert.Equal(t, err, io.EOF)
}

func TestCSVCustomDate(t *testing.T) {
	c, err := New("", &Options{
		DTIndex:  0,
		DTFormat: "2006-01-02",
		Reader:   strings.NewReader(CustomDate),
	})
	assert.NoError(t, err)
	values, err := c.Next()
	assert.NoError(t, err)
	assert.Equal(t, 1750, values[0].Time().Year())
	assert.Equal(t, time.January, values[0].Time().Month())
	assert.Equal(t, 1, values[0].Time().Day())
	values, _ = c.Next()
	assert.Equal(t, time.February, values[0].Time().Month())
	values, _ = c.Next()
	assert.Equal(t, time.March, values[0].Time().Month())
	values, _ = c.Next()
	assert.Equal(t, time.April, values[0].Time().Month())
	values, _ = c.Next()
	assert.Equal(t, time.May, values[0].Time().Month())
}
