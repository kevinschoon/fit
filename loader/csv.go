package loader

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

type CSVOptions struct {
	Path    string
	Reader  io.Reader // Optionally directly pass an io.Reader
	Parsers map[int]Parser
}

type CSVReader struct {
	reader  *csv.Reader
	file    *os.File
	columns []string
	Options *CSVOptions // CSV Options
}

func (c CSVReader) Next() ([]float64, error) {
	values := make([]float64, len(c.columns))
	records, err := c.reader.Read()
	if err != nil { // io.EOF indicates we are finished loading
		return nil, err
	}
	for i, record := range records {
		if parser, ok := c.Options.Parsers[i]; ok {
			if value, err := parser.Parse(record); err == nil {
				values[i] = value
				continue
			}
		}
		if value, err := strconv.ParseFloat(record, 64); err == nil {
			values[i] = value
		}
	}
	return values, nil
}

func (c CSVReader) Columns() []string { return c.columns }

func (c CSVReader) Close() error {
	return c.file.Close()
}

func NewCSV(opts *CSVOptions) (CSVReader, error) {
	c := CSVReader{
		Options: opts,
	}
	if opts.Reader == nil {
		file, err := os.Open(opts.Path)
		if err != nil {
			return c, err
		}
		c.file = file
		c.reader = csv.NewReader(bufio.NewReader(file))
	} else {
		c.reader = csv.NewReader(opts.Reader)
	}
	// Read the first record in the CSV to load column names
	record, err := c.reader.Read()
	if err != nil {
		return c, err
	}
	c.columns = make([]string, len(record))
	for i, name := range record {
		c.columns[i] = name
	}
	return c, nil
}
