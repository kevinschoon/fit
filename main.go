package main

import (
	"fmt"
	"os"
)

const (
	TakeoutDir string = "Takeout/Fit/Activities"
	StaticDir  string = "./www"
)

func main() {
	if err := InitDB(TakeoutDir); err != nil {
		fmt.Println("ERROR: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Loaded %d laps\n", len(database.Laps(Filter(NullFilter))))
	RunServer(":8000")
}
