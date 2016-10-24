package loader

import (
	"errors"
	"fmt"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/kevinschoon/fit/parser"
	"github.com/kevinschoon/fit/types"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

var ErrUnequalValues = errors.New("unequal value size")

type Rower interface {
	Row() ([]string, error)
	Dims() (int, int)
}

type Options struct {
	Name    string
	Path    string
	Enc     string
	Columns []string
	Sheet   string // Sheet name (XLS)
	Size    int64  // File Size (XLS)
	Parsers map[int]parser.Parser
}

// Rower returns a Rower based on the configured options
func (opts *Options) Rower(fp *os.File) (Rower, error) {
	var split []string
	if opts.Name == "" {
		split = strings.Split(opts.Path, "/")
		opts.Name = split[len(split)-1]
		if strings.Contains(opts.Name, ".") {
			opts.Name = strings.Split(opts.Name, ".")[0]
		}
	}
	if opts.Enc == "" {
		split = strings.Split(opts.Path, "/")
		split = strings.Split(split[len(split)-1], ".")
		opts.Enc = split[len(split)-1]
	}
	switch {
	case opts.Enc == "csv":
		csv, err := NewCSV(fp)
		if err != nil {
			return nil, err
		}
		if len(opts.Columns) == 0 {
			opts.Columns = csv.Columns
		}
		return csv, nil
	case opts.Enc == "xls" || opts.Enc == "xlsx":
		xls, err := NewXLS(fp, *opts)
		if err != nil {
			return nil, err
		}
		if len(opts.Columns) == 0 {
			return nil, fmt.Errorf("Specify at least one column")
		}
		return xls, nil
	}
	panic(fmt.Sprintf("unknown encoding: %d", opts.Enc))
}

func Matrix(rower Rower, parsers map[int]parser.Parser) (*mtx.Dense, error) {
	r, c := rower.Dims()
	mx := mtx.NewDense(r, c, nil)
	for j := 0; j < r; j++ {
		strs, err := rower.Row()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		row := make([]float64, c)
		for i, str := range strs {
			if parser, ok := parsers[i]; ok {
				if value, err := parser.Parse(str); err == nil {
					row[i] = value
					continue
				}
			}
			if value, err := strconv.ParseFloat(str, 64); err == nil {
				row[i] = value
				continue
			}
			row[i] = math.NaN()
		}
		mx.SetRow(j, row)
	}
	return mx, nil
}

func ReadPath(opts Options) (*types.Dataset, error) {
	var (
		rower Rower
		mx    *mtx.Dense
	)
	fp, err := os.Open(opts.Path)
	if err != nil {
		return nil, err
	}
	stats, err := fp.Stat()
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	opts.Size = stats.Size()
	rower, err = opts.Rower(fp)
	if err != nil {
		return nil, err
	}
	mx, err = Matrix(rower, opts.Parsers)
	if err != nil {
		return nil, err
	}
	return &types.Dataset{
		Name:    opts.Name,
		Columns: opts.Columns,
		Mtx:     mx,
	}, nil
}
