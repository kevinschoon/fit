package loader

import "io"

type MemReader struct {
	count   int
	columns []string
	values  [][]float64
}

func (m MemReader) Columns() []string { return m.columns }

func (m MemReader) Close() error { return nil }

func (m *MemReader) Next() ([]float64, error) {
	if m.count >= len(m.values) {
		return nil, io.EOF
	}
	values := m.values[m.count]
	m.count++
	return values, nil
}

func NewMemReader(columns []string, values [][]float64) *MemReader {
	return &MemReader{
		columns: columns,
		values:  values,
	}
}
