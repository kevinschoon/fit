package server

import (
	"bytes"
	"encoding/json"
	"github.com/kevinschoon/fit/clients"
	"github.com/kevinschoon/fit/types"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
)

var cleanup bool = false

var (
	datasetRaw []byte = []byte(`{"Name": "TestDataset", "Columns": ["V1", "V2"]}`)
	valuesRaw  []byte = []byte(`[1.0, 1.0, 2.0, 2.0, 3.0, 3.0, 4.0, 4.0, 5.0, 5.0]`)
)

type MockWriter struct {
	fp *os.File
}

func (writer MockWriter) Header() http.Header { return http.Header{} }

func (writer MockWriter) Write(data []byte) (int, error) {
	return writer.fp.Write(data)
}
func (writer MockWriter) WriteHeader(int) {}

type MockReader struct {
	*bytes.Reader
}

func (reader MockReader) Close() error { return nil }

func NewTestDB(t *testing.T) (types.Client, func()) {
	f, err := ioutil.TempFile("/tmp", "fit-test-")
	assert.NoError(t, err)
	db, err := clients.NewBoltClient(f.Name())
	assert.NoError(t, err)
	return db, func() {
		if cleanup {
			os.Remove(f.Name())
		}
	}
}

func NewMockWriter(t *testing.T) (MockWriter, func()) {
	fp, err := ioutil.TempFile("/tmp", "fit-test-writer-")
	assert.NoError(t, err)
	return MockWriter{fp: fp}, func() {
		fp.Close()
		if cleanup {
			os.Remove(fp.Name())
		}
	}
}

func TestDatasetAPI(t *testing.T) {
	db, cleanup := NewTestDB(t)
	defer cleanup()
	handler := Handler{db: db}
	reader := MockReader{bytes.NewReader(datasetRaw)}
	assert.NoError(t, handler.DatasetAPI(MockWriter{fp: nil}, &http.Request{
		Method: "POST",
		Body:   io.ReadCloser(reader),
	}))
	writer, cleanup := NewMockWriter(t)
	defer cleanup()
	assert.NoError(t, handler.DatasetAPI(writer, &http.Request{
		URL:    &url.URL{},
		Method: "GET",
	}))
	raw, err := ioutil.ReadFile(writer.fp.Name())
	assert.NoError(t, err)
	ds := []*types.Dataset{}
	assert.NoError(t, json.Unmarshal(raw, &ds))
	assert.Len(t, ds, 1)
	assert.Equal(t, "V1", ds[0].Columns[0])
	assert.Equal(t, "V2", ds[0].Columns[1])
	assert.Equal(t, 0, ds[0].Stats.Rows)
	assert.Equal(t, 0, ds[0].Stats.Columns)
	assert.NoError(t, handler.DatasetAPI(MockWriter{fp: nil}, &http.Request{
		Method: "DELETE",
		URL: &url.URL{
			RawQuery: "name=TestDataset",
		},
	}))
}
