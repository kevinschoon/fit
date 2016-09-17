package server

/*

import (
	"github.com/kevinschoon/gofit/models"
	"net/url"
	"time"
)

func StartEnd(u *url.URL) (start time.Time, end time.Time) {
	start, _ = time.Parse("2006-Jan-02", u.Query().Get("start"))
	end, _ = time.Parse("2006-Jan-02", u.Query().Get("end"))
	if end.IsZero() {
		end = time.Now()
	}
	return start, end
}

func Aggr(u *url.URL) (aggr models.Aggregation) {
	aggr = models.None
	switch u.Query().Get("aggr") {
	case "days":
		aggr = models.Days
	case "months":
		aggr = models.Months
	case "years":
		aggr = models.Years
	}
	return aggr
}

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
*/
