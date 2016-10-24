package loader

import (
	"encoding/csv"
	"io"
)

type CSV struct {
	Columns []string
	reader  *csv.Reader
	rows    [][]string
	index   int
}

func (c *CSV) Row() ([]string, error) {
	if c.index == len(c.rows) {
		return nil, io.EOF
	}
	row := c.rows[c.index]
	c.index++
	return row, nil
}

func (c CSV) Dims() (int, int) {
	return len(c.rows), len(c.Columns)
}

func NewCSV(reader io.Reader) (*CSV, error) {
	c := &CSV{
		reader: csv.NewReader(reader),
	}
	// Read the first record in the CSV to load column names
	row, err := c.reader.Read()
	if err != nil {
		return c, err
	}
	c.Columns = make([]string, len(row))
	for i, name := range row {
		c.Columns[i] = name
	}
	// Load the entire CSV into memory so
	// we can get it's demensions
	for {
		row, err := c.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		c.rows = append(c.rows, row)
	}
	return c, nil
}
