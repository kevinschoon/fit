package server

import (
	"fmt"
	"github.com/gonum/plot/vg"
	"github.com/kevinschoon/gofit/chart"
	"github.com/kevinschoon/gofit/models"
	"image/color"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// StartEnd extracts start and end times from a URL
// Default: 1970 - Now
func StartEnd(u *url.URL) (start time.Time, end time.Time) {
	start, _ = time.Parse(time.RFC3339, u.Query().Get("start"))
	end, _ = time.Parse(time.RFC3339, u.Query().Get("end"))
	if end.IsZero() {
		end = time.Now().UTC()
	}
	return start, end
}

// Aggr extracts an aggregation level from a URL
// Default: 24 hours
func Aggr(u *url.URL) (duration time.Duration) {
	raw := u.Query().Get("aggr")
	if raw == "none" {
		return time.Duration(0)
	}
	duration, err := time.ParseDuration(raw)
	if err != nil {
		duration = 24 * time.Hour
	}
	return duration
}

// GetBool extracts a boolean flag from a URL
func GetBool(name string, def bool, u *url.URL) bool {
	if v := u.Query().Get(name); v != "" {
		switch v {
		case "true":
			return true
		case "1":
			return true
		}
	}
	return def
}

// StrArray extracts an array of strings from a URL
func StrArray(name string, u *url.URL) []string {
	out := strings.Split(u.Query().Get(name), ",")
	if len(out) == 1 {
		if out[0] == "" {
			return []string{}
		}
	}
	return out
}

// Fn extracts a fuction from the URL
// Default: Avg
func Fn(u *url.URL) models.Function {
	switch u.Query().Get("fn") {
	case "sum":
		return models.Function(models.Sum)
	case "avg":
		return models.Function(models.Avg)
	case "min":
		return models.Function(models.Min)
	case "max":
		return models.Function(models.Max)
	default:
		return models.Function(models.Avg)
	}
}

// ChartCfg extracts a charting configuration from the URL
func ChartCfg(series *models.Series, u *url.URL) chart.Config {
	cfg := chart.Config{
		Title:          series.Name,
		PrimaryColor:   color.White,
		SecondaryColor: color.Black,
		Width:          18 * vg.Inch,
		Height:         5 * vg.Inch,
	}
	cfg.XAxis = models.Key(0) // Default
	name := u.Query().Get("X")
	if key, ok := series.Keys[name]; ok {
		cfg.XLabel = name
		cfg.XAxis = key
	}
	if name == "time" {
		cfg.PlotTime = true
	}
	cfg.YAxis = make(map[string]models.Key)
	for _, name := range StrArray("Y", u) {
		if key, ok := series.Keys[name]; ok {
			cfg.YAxis[name] = key
		}
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
		Path:     fmt.Sprintf("%s/chart", u.Path),
		RawQuery: u.Query().Encode(),
	}
	return c.String()
}
