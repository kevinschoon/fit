package types

import (
	"encoding/json"
	"fmt"
	mtx "github.com/gonum/matrix/mat64"
	"strconv"
	"strings"
	"time"
)

type grouping struct {
	Name  string
	Index int
	Max   string
}

// Grouping represents a "group by" configuration
type Grouping struct {
	Name  string
	Index int
	Max   time.Duration
}

func (grp Grouping) UnmarshalJSON(data []byte) error {
	in := &grouping{}
	if err := json.Unmarshal(data, in); err != nil {
		return err
	}
	grp.Name = in.Name
	grp.Index = in.Index
	max, _ := time.ParseDuration(in.Max)
	grp.Max = max
	return nil
}

func (grp Grouping) MarshalJSON() ([]byte, error) {
	return json.Marshal(&grouping{
		Name:  grp.Name,
		Index: grp.Index,
		Max:   grp.Max.String(),
	})
}

func (grp Grouping) String() string {
	return fmt.Sprintf("%s,%d,%s", grp.Name, grp.Index, grp.Max.String())
}

func (grp Grouping) Group(other *mtx.Dense) []mtx.Matrix {
	r, c := other.Dims()
	views := make([]mtx.Matrix, 0)
	var (
		view     mtx.Matrix
		previous time.Time
		current  time.Time
		duration time.Duration
	)
	for i, j := 0, 1; i+j <= r; j++ {
		view = other.View(i, 0, j, c)
		current = time.Unix(int64(view.At(j-1, grp.Index)), 0).UTC()
		duration += current.Sub(previous)
		if duration >= grp.Max {
			views = append(views, view)
			i, j = i+j, 0
			duration = time.Duration(0)
		}
		previous = current
	}
	return views
}

// NewGrouping returns a grouping based on a string
// parameter. See documentation for Duration str format
//
// Duration,0,1min
// ^--------^-^----Name,Index,DurationStr
func NewGrouping(arg string) *Grouping {
	split := strings.Split(arg, ",")
	var grouping *Grouping
	if len(split) > 0 {
		grouping = &Grouping{
			Name: split[0],
		}
	}
	if len(split) >= 2 {
		index, _ := strconv.ParseInt(split[1], 0, 64)
		grouping.Index = int(index)
	}
	if len(split) >= 3 {
		duration, _ := time.ParseDuration(split[2])
		grouping.Max = duration
	}
	return grouping
}
