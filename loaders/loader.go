package loaders

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/kevinschoon/fit/loaders/csv"
	"github.com/kevinschoon/fit/models"
	"io"
	"os"
	"strings"
)

type FileType int

const (
	None FileType = iota
	CSV
	TCX
)

func FileTypeByName(name string) FileType {
	switch name {
	case "csv":
		return CSV
	case "tcx":
		return TCX
	}
	return None
}

var ErrUnknownFileType = errors.New("Unknown File Type")

type Options struct {
	Name       string       // Name of the Series/Dataset
	Path       string       // Path for file or directory
	Type       FileType     // Optional file type
	CSVOptions *csv.Options // CSV-specific options
	Stdout     bool         // Dump to stdout
	Values     bool         // Dump values instead of just Series
}

type Loader interface {
	Next() ([]models.Value, error)
	Keys() models.Keys
	Close() error
}

func Load(opts *Options) (Loader, error) {
	// Attempt to detect FileType if not provided
	if opts.Type == None {
		split := strings.Split(opts.Path, ".")
		opts.Type = FileTypeByName(split[len(split)-1])
	}
	// Attempt to detect name for series
	if opts.Name == "" {
		split := strings.Split(opts.Path, "/")
		split = strings.Split(split[len(split)-1], ".")
		opts.Name = split[0]
	}
	switch opts.Type {
	case CSV:
		return csv.New(opts.Path, opts.CSVOptions)
	}
	return nil, ErrUnknownFileType
}

// Stdout iterates a Loader object until EOF
// If values is false it will dump an empty Series object
// if values is true it will dump all values of a series
func Stdout(name string, values bool, loader Loader) error {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", " ")
	if values {
		for {
			values, err := loader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if err := encoder.Encode(values); err != nil {
				return err
			}
		}
		return nil
	}
	series := models.NewSeries(loader.Keys().Names())
	series.Name = name
	return encoder.Encode(series)
}
