package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"testing"
)

// Obj is a test object that implements all required DB interfaces
type Obj struct {
	Name string
	X    float64
	Y    float64
}

type Objs struct {
	objs []Obj
}

func (objs Objs) Keys() []string {
	return []string{"Name", "X", "Y"}
}

func (objs Objs) Values() [][]string {
	values := make([][]string, len(objs.objs))
	for i, obj := range objs.objs {
		values[i] = make([]string, len(objs.Keys()))
		values[i][0] = obj.Name
		values[i][1] = fmt.Sprintf("%d", int(obj.X))
		values[i][2] = fmt.Sprintf("%d", int(obj.Y))
	}
	return values
}

func (objs Objs) Pts(key string) Datapoints {
	pts := make(Datapoints, len(objs.objs))
	for i, obj := range objs.objs {
		pts[i].X = obj.X
		pts[i].Y = obj.Y
	}
	return pts
}

func (objs Objs) Write(db *gorm.DB) error {
	for _, obj := range objs.objs {
		if err := db.Create(&obj).Error; err != nil {
			return err
		}
	}
	return nil
}

func (objs Objs) Read(db *gorm.DB) (Series, error) {
	if err := db.Find(&objs.objs).Error; err != nil {
		return nil, err
	}
	return objs, nil
}

func (objs Objs) Query(Query) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("x > ?", 0.4)
	}
}

func newDB(t *testing.T, values ...interface{}) *gorm.DB {
	f, err := ioutil.TempFile("/tmp", "gofit")
	if err != nil {
		t.Error(err)
	}
	f.Close()
	db, err := New(f.Name(), values...)
	if err != nil {
		t.Error(err)
	}
	return db
}

func TestDB(t *testing.T) {
	db := newDB(t, &Obj{})
	objs := Objs{objs: []Obj{Obj{Name: "Hello", X: 0.5, Y: 1.5}}}
	if err := Write(db, objs); err != nil {
		t.Error(err)
	}
	series, err := Read(db, Query{}, Objs{})
	if err != nil {
		t.Error(err)
	}
	if len(series.Keys()) != 3 {
		t.Errorf("Keys should be 3")
	}
	if len(series.Values()) != 1 {
		t.Errorf("values should equal 1")
	}
	if len(series.Pts("")) != 1 {
		t.Errorf("datapoints should equal 1")
	}
}
