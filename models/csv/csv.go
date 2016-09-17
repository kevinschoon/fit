/*
Loads arbitrary CSV data into a Series
*/
package csv

import (
	"encoding/csv"
	"github.com/kevinschoon/gofit/models"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

// CSV stores records from CSV files
type CSV struct {
	records [][]string
}

// Load returns an array of Series from loaded CSV data
func (c *CSV) Load() []*models.Series {
	series := make([]*models.Series, 1)
	names := make([]string, 0)
	start := time.Now()
	for _, name := range c.records[0] {
		if name == "" {
			name = "NO_NAME"
		}
		names = append(names, name)
	}
	series[0] = models.NewSeries(names)
	for _, record := range c.records[1:] {
		values := make([]models.Value, 0)
		for _, v := range record {
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				continue
			}
			// TODO: This will panic if the length of
			// values changes. Need to support missing data
			values = append(values, models.Value(value))
		}
		series[0].Add(start, values)
	}
	return series
}

// FromFile loads CSV data from a single file and returns a CSV
func FromFile(path string) (*CSV, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var records [][]string
	reader := csv.NewReader(strings.NewReader(string(raw)))
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return &CSV{records: records}, nil
}

// FromDir loads discovers CSV files in a directory and returns a CSV
func FromDir(path string) (*CSV, error) {
	// TODO
	return nil, nil
}
