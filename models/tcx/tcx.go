package tcx

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kevinschoon/gofit/models"
	"github.com/kevinschoon/tcx"
	"sort"
)

// Data satisfies the "series" model interface
type Data struct {
	acts      tcx.Acts
	precision models.Precision
}

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

// Loader satisfy database interfaces
type Loader struct {
	Path   string
	tcxdbs []*tcx.TCXDB
}

func (t Loader) Types() []interface{} {
	return []interface{}{&tcx.Activity{}, &tcx.Lap{}, &tcx.Track{}, &tcx.Trackpoint{}}
}

func (t Loader) Write(db *gorm.DB) error {
	for _, tcxdb := range t.tcxdbs {
		for _, activity := range tcxdb.Acts.Act {
			if err := db.Create(&activity).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (t Loader) query(query models.Query) func(*gorm.DB) *gorm.DB {
	qs := "DATE(start_time) >= ? AND DATE(start_time) <= ?"
	values := []interface{}{
		query.Start.Format("2006-01-02"),
		query.End.Format("2006-01-02"),
	}
	if len(query.Match) == 1 {
		for key, value := range query.Match {
			switch key {
			case "activity":
				qs += fmt.Sprintf(" AND %s LIKE ?", key)
				values = append(values, value)
			default:
				break // Column name must be verified since it cannot be escaped
			}
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Laps.Trk.Pt").Where(qs, values...)
	}
}

func (t Loader) Read(db *gorm.DB, query models.Query) (models.Series, error) {
	var (
		count int
		last  int
	)
	data := Data{precision: query.Precision}
	fn := t.query(query)
	fn(db).Model(&tcx.Activity{}).Count(&count)
	for len(data.acts) < count {
		if db.Error != nil {
			return data, db.Error
		}
		var results tcx.Acts
		fn(db).Limit(100).Offset(last).Find(&results)
		for _, result := range results {
			data.acts = append(data.acts, result)
		}
		last += 100
	}
	return data, nil
}

// FromDir loads TCX data from the given directory
func (t *Loader) FromDir(path string) error {
	tcxDbs, err := tcx.ReadDir(path)
	if err != nil {
		return err
	}
	for _, db := range tcxDbs {
		t.tcxdbs = append(t.tcxdbs, db)
	}
	return nil
}
