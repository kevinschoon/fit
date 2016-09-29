package csv

import (
	"encoding/csv"
	"github.com/kevinschoon/fit/parser"
	"io"
	"strconv"
)

type CSV struct {
	reader  *csv.Reader
	columns []string
	parsers map[int]parser.Parser
}

func (c CSV) Next() ([]float64, error) {
	values := make([]float64, len(c.columns))
	records, err := c.reader.Read()
	if err != nil { // io.EOF indicates we are finished loading
		return nil, err
	}
	for i, record := range records {
		if parser, ok := c.parsers[i]; ok {
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

func (c CSV) Columns() []string { return c.columns }

func New(reader io.Reader, parsers map[int]parser.Parser) (CSV, error) {
	c := CSV{
		parsers: parsers,
		reader:  csv.NewReader(reader),
	}
	if c.parsers == nil {
		c.parsers = make(map[int]parser.Parser)
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
