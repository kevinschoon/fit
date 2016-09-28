package loader

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

func TestCSVReader(t *testing.T) {
	c, err := NewCSV(&CSVOptions{
		Reader: strings.NewReader(Simple),
	})
	defer c.Close()
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

func TestCSVReaderTimeParser(t *testing.T) {
	c, err := NewCSV(&CSVOptions{
		Parsers: map[int]Parser{
			0: TimeParser{
				Format: "2006-01-02",
			},
		},
		Reader: strings.NewReader(CustomDate),
	})
	assert.NoError(t, err)
	values, err := c.Next()
	assert.NoError(t, err)
	dt := time.Unix(int64(values[0]), 0).UTC()
	assert.Equal(t, 1750, dt.Year())
	assert.Equal(t, time.January, dt.Month())
	assert.Equal(t, 1, dt.Day())
}
