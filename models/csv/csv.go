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

type Options struct{}

type CSV struct {
	records [][]string
}

func (c *CSV) Load() *models.Collection {
	names := c.records[0]
	collection := &models.Collection{}
	for _, record := range c.records[1:] {
		values := make([]models.Value, 0)
		for i, v := range record {
			value, err := strconv.ParseFloat(v, 64)
			if err != nil {
				continue
			}
			// TODO: This will panic if the length of
			// values changes. Need to support missing data
			values = append(values, models.Value{
				Name:  names[i],
				Value: value,
			})
		}
		collection.Add(time.Now(), values)
	}
	return collection
}

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

func FromDir(path string) (*CSV, error) {
	// TODO
	return nil, nil
}
