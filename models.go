package main

import (
	"github.com/kevinschoon/tcx"
	"sort"
	"time"
)

type Datapoint struct {
	X float64
	Y float64
}

type Datapoints []Datapoint

func (dps Datapoints) Len() int {
	return len(dps)
}

func (dps Datapoints) Less(i, j int) bool {
	return dps[i].X < dps[j].X
}

func (dps Datapoints) Swap(i, j int) {
	dps[i], dps[j] = dps[j], dps[i]
}

func (dps Datapoints) Value(i int) float64 {
	return dps[i].Y
}

func (dps Datapoints) XY(i int) (x, y float64) {
	return dps[i].X, dps[i].Y
}

func ActivityByDist(acts tcx.Acts) Datapoints {
	dps := make(Datapoints, len(acts))
	for i, act := range acts {
		dps[i] = Datapoint{X: float64(act.StartTime.Unix())}
		for _, lap := range act.Laps {
			dps[i].Y += lap.Dist
		}
	}
	//sort.Sort(dps)
	return dps
}

func RollUpActivities(acts tcx.Acts, precision string) tcx.Acts {
	if precision == "none" {
		return acts
	}
	buckets := make(map[int][]tcx.Activity)
	for _, act := range acts {
		key := TimeKey(act.StartTime, precision)
		if _, ok := buckets[key]; !ok {
			buckets[key] = make([]tcx.Activity, 0)
		}
		buckets[key] = append(buckets[key], act)
	}
	activities := tcx.Acts{}
	for _, bucket := range buckets {
		first := bucket[0]
		if len(bucket) > 1 {
			bucket = bucket[1:]
			for _, act := range bucket {
				for _, lap := range act.Laps {
					first.Laps = append(first.Laps, lap)
				}
			}
		}
		activities = append(activities, first)
	}
	sort.Sort(activities)
	return activities
}

// Generate a unique key for a given date
func TimeKey(t time.Time, precision string) int {
	var key int
	switch precision {
	case "year":
		key = t.Year()
	case "month":
		key = (t.Year() * 12) + int(t.Month())
	case "day":
		key = (t.Year() * 12) + (int(t.Month()) * 31) + t.Day()
	}
	return key
}
