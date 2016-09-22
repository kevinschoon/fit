/*
Loads arbitrary CSV data into a Series
*/
package csv

import (
	"bufio"
	"encoding/csv"
	"github.com/kevinschoon/fit/models"
	"os"
	"strconv"
	"time"
)

type Options struct {
	DTFormat string // DateTime Format
	DTIndex  int    // Record index to parse the date from
}

// CSVLoader implements the Loader interface
type CSVLoader struct {
	count   int
	reader  *csv.Reader
	file    *os.File
	keys    models.Keys
	Options *Options // CSV Options
}

func (c CSVLoader) Next() ([]models.Value, error) {
	records, err := c.reader.Read()
	if err != nil { // io.EOF indicates we are finished loading
		return nil, err
	}
	var values []models.Value
	for i, record := range records {
		if i == c.Options.DTIndex && c.Options.DTFormat != "" {
			start, err := time.Parse(c.Options.DTFormat, record)
			if err != nil {
				return nil, err
			}
			values = append(values, models.Value(start.Unix()))
			continue
		}
		value, _ := strconv.ParseFloat(record, 64)
		values = append(values, models.Value(value))
	}
	return values, nil
}

func (c CSVLoader) Keys() models.Keys { return c.keys }

func (c CSVLoader) Close() error {
	return c.file.Close()
}

// New creates a new CSVLoader for the path provided
func New(path string, opts *Options) (CSVLoader, error) {
	loader := CSVLoader{
		keys:    make(models.Keys),
		Options: opts,
	}
	file, err := os.Open(path)
	if err != nil {
		return loader, err
	}
	loader.file = file
	loader.reader = csv.NewReader(bufio.NewReader(file))
	// Read the first record in the CSV to load column names
	record, err := loader.reader.Read()
	if err != nil {
		return loader, err
	}
	for i, name := range record {
		if name == "" {
			name = "NO_NAME"
		}
		loader.keys[name] = models.Key(i)
	}
	return loader, err
}
