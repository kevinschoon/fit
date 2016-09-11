package tcx

import (
	"github.com/kevinschoon/gofit/models"
	"github.com/kevinschoon/tcx"
)

const (
	Laps models.Key = iota
	Distance
	Duration
)

type TCX struct {
	Acts tcx.Acts
}

func (t TCX) Load() *models.Collection {
	collection := &models.Collection{}
	for _, act := range t.Acts {
		collection.Add(act.StartTime, []models.Value{
			models.Value{
				Name:  "Laps",
				Value: float64(len(act.Laps)),
			},
			models.Value{
				Name:  "Dist",
				Value: act.Distance(),
			},
			models.Value{
				Name:  "Duration",
				Value: act.Duration(),
			},
		})
	}
	return collection
}

// FromDir loads TCX data from a directory
func FromDir(path string) (*TCX, error) {
	tcxDbs, err := tcx.ReadDir(path)
	if err != nil {
		return nil, err
	}
	t := &TCX{}
	for _, db := range tcxDbs {
		for _, act := range db.Acts.Act {
			t.Acts = append(t.Acts, act)
		}
	}
	return t, nil
}
