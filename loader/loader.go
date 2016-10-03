package loader

import (
	"errors"
	"fmt"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/kevinschoon/fit/loader/csv"
	"github.com/kevinschoon/fit/parser"
	"github.com/kevinschoon/fit/types"
	"io"
	"os"
	"strings"
)

var ErrUnequalValues = errors.New("unequal value size")

type Encoding int

const (
	NONE Encoding = iota
	CSV
	TCX
	BYTES
	JSON
)

// Loader provides an iterative interface
// for loading pairs of float64 values
type Loader interface {
	Next() ([]float64, error)
	Columns() []string
}

func Load(loader Loader) (*mtx.Dense, error) {
	var values []float64
	var rows int
	width := len(loader.Columns())
	for {
		v, err := loader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(v) != width {
			return nil, ErrUnequalValues
		}
		for _, value := range v {
			values = append(values, value)
		}
		rows++
	}
	return mtx.NewDense(rows, width, values), nil
}

func ReadPath(name, path string, enc Encoding, parsers map[int]parser.Parser) (*types.Dataset, error) {
	var (
		loader Loader
		mx     *mtx.Dense
	)
	if name == "" {
		split := strings.Split(path, "/")
		name = split[len(split)-1]
	}
	if enc == NONE {
		split := strings.Split(path, ".")
		switch split[len(split)-1] {
		case "csv":
			enc = CSV
		}
	}
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	switch enc {
	case CSV:
		loader, err = csv.New(fp, parsers)
		if err != nil {
			return nil, err
		}
		mx, err = Load(loader)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("unknown encoding: %d", enc))
	}
	return &types.Dataset{
		Name:    name,
		Columns: loader.Columns(),
		Mtx:     mx,
	}, nil
}
