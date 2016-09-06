package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/kevinschoon/gofit/models"
	"github.com/kevinschoon/tcx"
	"sort"
)

type TCXData struct {
	acts      tcx.Acts
	precision models.Precision
}

func (d TCXData) Columns() []string {
	return []string{"Laps", "Distance", "Duration"}
}

func (d TCXData) Rows() models.Rows {
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

func (d TCXData) Pts(key string) models.Datapoints {
	rows := d.Rows()
	sort.Sort(rows)
	pts := make(models.Datapoints, len(rows))
	for i, row := range rows {
		pts[i].X = float64(row.Time.Unix())
		pts[i].Y = float64(row.Values[1])
	}
	return pts
}

type TCXLoader struct {
	Path   string
	tcxdbs []*tcx.TCXDB
}

func (t TCXLoader) Write(db *gorm.DB) error {
	for _, tcxdb := range t.tcxdbs {
		for _, activity := range tcxdb.Acts.Act {
			if err := db.Create(&activity).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (t TCXLoader) query(query Query) func(*gorm.DB) *gorm.DB {
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

func (t TCXLoader) Read(db *gorm.DB, query Query) (models.Series, error) {
	var (
		count int
		last  int
	)
	data := TCXData{precision: query.Precision}
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

func (t *TCXLoader) Load() error {
	tcxDbs, err := tcx.ReadDir(t.Path)
	if err != nil {
		return err
	}
	for _, db := range tcxDbs {
		t.tcxdbs = append(t.tcxdbs, db)
	}
	return nil
}
