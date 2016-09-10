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

/*
// Write implements database.Writer
func (t TCX) Write(db *gorm.DB) error {
	for _, act := range t.acts {
		if err := db.Create(&act).Error; err != nil {
			return err
		}
	}
	return nil
}

// Read implements database.Reader
func (t TCX) Read(db *gorm.DB, query models.Query) (models.Serieser, error) {
	var (
		count int
		last  int
	)
	//data := Data{precision: query.Precision}
	fn := getQuery(query)
	fn(db).Model(&tcx.Activity{}).Count(&count)
	for len(t.acts) < count {
		if db.Error != nil {
			return t, db.Error
		}
		var results tcx.Acts
		fn(db).Limit(100).Offset(last).Find(&results)
		for _, result := range results {
			t.acts = append(t.acts, result)
		}
		last += 100
	}
	return t, db.Error
}
*/

/*
// Types implements models.Serieser
func (t TCX) Types() []interface{} {
	return []interface{}{&tcx.Activity{}, &tcx.Lap{}, &tcx.Track{}, &tcx.Trackpoint{}}
}

// getQuery returns a gorm-compatible query
func getQuery(query models.Query) func(*gorm.DB) *gorm.DB {
	qs := "DATE(start_time) >= ? AND DATE(start_time) <= ?"
	values := []interface{}{
		query.Start.Format("2006-01-02"),
		query.End.Format("2006-01-02"),
	}
	for key, value := range query.Match {
		switch key {
		case "sport":
			qs += fmt.Sprintf(" AND %s LIKE ?", key)
			values = append(values, value)
		default:
			break // Column name must be verified since it cannot be escaped
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Laps.Trk.Pt").Where(qs, values...)
	}
}
*/
/*
func (d Data) Columns() []string {
	return []string{"Laps", "Distance", "Duration"}
}

func (d Data) Rows() models.Rows {
	values := make(models.Rows, len(d.acts))
	for i, act := range d.acts {
		values[i] = models.Row{
			Time:   act.StartTime,
			Values: make([]models.Value, 3),
		}
		values[i].Values[0] = models.Value(len(act.Laps))
		values[i].Values[1] = models.Value(act.Distance())
		values[i].Values[2] = models.Value(act.Duration())
	}
	values = values.RollUp(d.precision)
	sort.Sort(sort.Reverse(values))
	return values
}

func (d Data) Pts(key string) models.Datapoints {
	rows := d.Rows()
	sort.Sort(rows)
	pts := make(models.Datapoints, len(rows))
	for i, row := range rows {
		pts[i].X = float64(row.Time.Unix())
		pts[i].Y = float64(row.Values[1])
	}
	return pts
}
*/
