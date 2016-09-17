package server

import (
	"github.com/kevinschoon/gofit/models"
	"net/url"
	"strings"
	"time"
)

// StartEnd extracts start and end times from a URL
// Default: 1970 - Now
func StartEnd(u *url.URL) (start time.Time, end time.Time) {
	start, _ = time.Parse("2006-Jan-02", u.Query().Get("start"))
	end, _ = time.Parse("2006-Jan-02", u.Query().Get("end"))
	if end.IsZero() {
		end = time.Now()
	}
	return start, end
}

// Aggr extracts an aggregation level from a URL
// Default: 24 hours
func Aggr(u *url.URL) (duration time.Duration) {
	duration, err := time.ParseDuration(u.Query().Get("aggr"))
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

/*
func XY(u *url.URL, collection *models.Collection) (X models.Key, Y models.Key) {
	X = models.Key(0)
	Y = models.Key(1)
	if name := u.Query().Get("x"); name != "" {
		for i, n := range collection.Names() {
			if name == n {
				X = models.Key(i)
			}
		}
	}
	if name := u.Query().Get("y"); name != "" {
		for i, n := range collection.Names() {
			if name == n {
				Y = models.Key(i)
			}
		}
	}
	return X, Y
}
*/
