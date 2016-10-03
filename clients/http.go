package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kevinschoon/fit/types"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const timeout time.Duration = 10 * time.Second

// HTTPClient implements the types.Client
// interface over HTTP
type HTTPClient struct {
	baseURL *url.URL
	client  *http.Client
}

func (c *HTTPClient) url() *url.URL {
	return &url.URL{
		Scheme: c.baseURL.Scheme,
		Host:   c.baseURL.Host,
		Path:   "/1/dataset",
	}
}

func (c *HTTPClient) do(req *http.Request) ([]byte, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		log.Println("API Error: %d - %s", res.StatusCode, data)
		switch res.StatusCode {
		case 404:
			return nil, types.ErrNotFound
		default:
			return nil, types.ErrAPI
		}
	}
	return data, nil
}

func (c *HTTPClient) Datasets() (datasets []*types.Dataset, err error) {
	raw, err := c.do(&http.Request{
		URL:    c.url(),
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(raw, &datasets)
	return datasets, err
}

func (c *HTTPClient) Write(ds *types.Dataset) (err error) {
	ds.WithValues = true
	raw, err := json.Marshal(ds)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.url().String(), io.Reader(bytes.NewReader(raw)))
	if err != nil {
		return err
	}
	_, err = c.do(req)
	return err
}

func (c *HTTPClient) Delete(name string) (err error) {
	u := c.url()
	u.RawQuery = fmt.Sprintf("name=%s", name)
	_, err = c.do(&http.Request{
		URL:    u,
		Method: "DELETE",
	})
	return err
}

func (c *HTTPClient) Query(queries types.Queries) (ds *types.Dataset, err error) {
	u := c.url()
	u.RawQuery = queries.QueryStr()
	raw, err := c.do(&http.Request{
		URL:    u,
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}
	ds = &types.Dataset{
		WithValues: true,
	}
	err = json.Unmarshal(raw, ds)
	return ds, err
}

func NewHTTPClient(endpoint string) (types.Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	c := &HTTPClient{
		baseURL: u,
		client: &http.Client{
			Timeout: timeout,
		},
	}
	return types.Client(c), nil
}
