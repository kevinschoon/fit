package loaders

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kevinschoon/gofit/loaders/csv"
	"github.com/kevinschoon/gofit/models"
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
	Close() func() error
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
	fp, err := os.Open(opts.Path)
	if err != nil {
		return nil, err
	}
	closer := func() error {
		if err := fp.Close(); err != nil {
			fmt.Println("Error: ", err)
			return err
		}
		return nil
	}
	reader := bufio.NewReader(fp)
	switch opts.Type {
	case CSV:
		return csv.New(reader, closer, opts.CSVOptions)
	}
	return nil, ErrUnknownFileType
}

func ToStdout(opts *Options, loader Loader) error {
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", " ")
	if opts.Values {
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
	series.Name = opts.Name
	return encoder.Encode(series)
}
