package server

import (
	"github.com/gonum/plot/vg"
	"github.com/kevinschoon/fit/chart"
	"github.com/kevinschoon/fit/store"
	"image/color"
	"net/url"
	"strconv"
)

func Queries(u *url.URL, name string) []string {
	values := u.Query()
	if v, ok := values[name]; ok {
		return v
	}
	return []string{}
}

func XYQueries(u *url.URL) store.Queries {
	values := u.Query()
	args := make([]string, 0)
	args = append(args, values.Get("x"))
	if v, ok := values["y"]; ok {
		for _, arg := range v {
			args = append(args, arg)
		}
	}
	return store.NewQueries(args)
}

// ChartCfg returns a chart.Config object from the URL
func ChartCfg(ds *store.Dataset, u *url.URL) chart.Config {
	cfg := chart.Config{
		Title:          ds.Name,
		PrimaryColor:   color.White,
		SecondaryColor: color.Black,
		Width:          18 * vg.Inch,
		Height:         5 * vg.Inch,
		Columns:        XYQueries(u).Columns(),
	}

	if w, err := strconv.ParseInt(u.Query().Get("width"), 0, 64); err == nil {
		if w < 20 { // Prevent potentially horrible DOS
			cfg.Width = vg.Length(w) * vg.Inch
		}
	}
	if h, err := strconv.ParseInt(u.Query().Get("height"), 0, 64); err == nil {
		if h < 20 {
			cfg.Height = vg.Length(h) * vg.Inch
		}
	}
	return cfg
}

// Chart builds the URL for rendering the chart
func Chart(u *url.URL) string {
	c := &url.URL{
		Scheme:   u.Scheme,
		Host:     u.Host,
		Path:     "/chart",
		RawQuery: u.Query().Encode(),
	}
	return c.String()
}
