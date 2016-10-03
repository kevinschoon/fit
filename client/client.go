package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kevinschoon/fit/store"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const timeout time.Duration = 10 * time.Second

var ErrAPIError = errors.New("api error")

// Client talks to the Fit REST API
// It implements the same interface as
// the store.DB. As the API stabilizes
// both should be factored into proper
// interface types.
type Client struct {
	baseURL *url.URL
	client  *http.Client
}

func (c *Client) url() *url.URL {
	return &url.URL{
		Scheme: c.baseURL.Scheme,
		Host:   c.baseURL.Host,
		Path:   "/1/dataset",
	}
}

func (c *Client) do(req *http.Request) ([]byte, error) {
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
			return nil, store.ErrNotFound
		default:
			return nil, ErrAPIError
		}
	}
	return data, nil
}

func (c *Client) Datasets() (datasets []*store.Dataset, err error) {
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

func (c *Client) Write(ds *store.Dataset) (err error) {
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

func (c *Client) Read(name string) (ds *store.Dataset, err error) {
	return nil, nil
}

func (c *Client) Delete(name string) (err error) {
	u := c.url()
	u.RawQuery = fmt.Sprintf("name=%s", name)
	_, err = c.do(&http.Request{
		URL:    u,
		Method: "DELETE",
	})
	return err
}

func (c *Client) Query(queries store.Queries) (ds *store.Dataset, err error) {
	u := c.url()
	u.RawQuery = queries.QueryStr()
	raw, err := c.do(&http.Request{
		URL:    u,
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}
	ds = &store.Dataset{
		WithValues: true,
	}
	err = json.Unmarshal(raw, ds)
	return ds, err
}

func NewClient(endpoint string) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	client := &Client{
		baseURL: u,
		client: &http.Client{
			Timeout: timeout,
		},
	}
	return client, nil
}
