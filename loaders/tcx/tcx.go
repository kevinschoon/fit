/*
Reads Garmin TCX data into a Series
*/
package tcx

import (
	"github.com/kevinschoon/gofit/models"
	"github.com/kevinschoon/tcx"
)

// TCX stores TCX activity data
type TCX struct {
	Name string
	Acts tcx.Acts
}

// Load returns an array of Series from loaded TCX data
func (t TCX) Load() []*models.Series {
	series := make([]*models.Series, 1)
	series[0] = models.NewSeries([]string{
		"Laps",
		"Distance",
		"Duration",
	})
	series[0].Name = t.Name
	for _, act := range t.Acts {
		series[0].Add(act.StartTime, []models.Value{
			models.Value(len(act.Laps)),
			models.Value(act.Distance()),
			models.Value(act.Duration()),
		})
	}
	return series
}

// FromDir loads TCX data from a directory
func FromDir(path, name string) (*TCX, error) {
	tcxDbs, err := tcx.ReadDir(path)
	if err != nil {
		return nil, err
	}
	t := &TCX{
		Name: name,
	}
	for _, db := range tcxDbs {
		for _, act := range db.Acts.Act {
			t.Acts = append(t.Acts, act)
		}
	}
	return t, nil
}
