package database

import (
	"fmt"
	"github.com/kevinschoon/tcx"
	"testing"
	"time"
)

func TestTCXLoader(t *testing.T) {
	db := newDB(t, tcx.Activity{}, tcx.Lap{}, tcx.Track{}, tcx.Trackpoint{})
	loader := &TCXLoader{Path: "./test/sample.tcx"}
	if err := loader.Load(); err != nil {
		t.Error(err)
	}
	fmt.Println(loader.tcxdbs)
	if err := Write(db, loader); err != nil {
		t.Error(err)
	}
	series, err := Read(db, Query{
		Start:  time.Now().AddDate(-100, 0, 0),
		End:    time.Now(),
		Column: "sport",
		Value:  "Running",
	}, loader)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(series.Values())
	for _, pt := range series.Pts("") {
		fmt.Println(pt.X, pt.Y)
	}
}
