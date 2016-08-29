package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/kevinschoon/tcx"
	"os"
)

func FailOnErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}

func Server(cmd *cli.Cmd) {
	pattern := cmd.StringOpt("pattern", ":8000", "IP and port pattern to listen on")
	cmd.Action = func() {
		RunServer(*pattern)
	}
}

func Load(cmd *cli.Cmd) {
	path := cmd.StringOpt("t tcxPath", "Takeout/Fit/Activities", "Path to TCX Data")
	cmd.Action = func() {
		database, err := GetDB()
		defer database.Close()
		FailOnErr(err)
		dbs, err := tcx.ReadDir(*path)
		FailOnErr(err)
		for _, db := range dbs {
			FailOnErr(BulkUpsert(database, db.Acts.Act))
		}
	}
}

func main() {
	app := cli.App("gofit", "GoFit!")
	var (
		Debug  = app.BoolOpt("d debug", true, "Debug Mode")
		DBPath = app.StringOpt("p path", "/tmp/gofit.db", "Path to SQLite DB")
	)
	options = &Options{
		Debug:  Debug,
		DBPath: DBPath,
	}
	app.Command("server", "Run the GoFit web UI", Server)
	app.Command("load", "Load TCX Data", Load)
	InitDB()
	app.Run(os.Args)
}
