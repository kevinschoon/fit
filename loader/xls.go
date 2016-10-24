package loader

import (
	"errors"
	"github.com/tealeg/xlsx"
	"io"
)

var InvalidXLS = errors.New("Invalid XLS")

type XLS struct {
	file    *xlsx.File
	sheet   *xlsx.Sheet
	index   int
	columns []string
}

func (x *XLS) Row() ([]string, error) {
	if x.index > x.sheet.MaxRow {
		return nil, io.EOF
	}
	values := make([]string, len(x.columns))
	for i := 0; i < len(values); i++ {
		if str, err := x.sheet.Cell(x.index, i).String(); err == nil {
			values[i] = str
		}
	}
	x.index++
	return values, nil
}

func (x XLS) Dims() (int, int) {
	return x.sheet.MaxRow, len(x.columns)
}

func NewXLS(reader io.ReaderAt, opts Options) (*XLS, error) {
	f, err := xlsx.OpenReaderAt(reader, opts.Size)
	if err != nil {
		return nil, err
	}
	xls := &XLS{
		columns: opts.Columns,
	}
	if opts.Sheet != "" {
		if sheet, ok := f.Sheet[opts.Sheet]; ok {
			xls.sheet = sheet
		} else {
			return nil, InvalidXLS
		}
	} else {
		if len(f.Sheets) == 0 {
			return nil, InvalidXLS
		}
		xls.sheet = f.Sheets[0]
	}
	return xls, nil
}
