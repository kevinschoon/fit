package cmd

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/kevinschoon/gofit/database"
	"github.com/kevinschoon/gofit/models/csv"
	"github.com/kevinschoon/gofit/models/tcx"
	"github.com/kevinschoon/gofit/server"
	"os"
)

func FailOnErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}

// Server starts the HTTP server
func Server(cmd *cli.Cmd) {
	var (
		pattern = cmd.StringOpt("pattern", ":8000", "IP and port pattern to listen on")
		path    = cmd.StringOpt("p path", "/tmp/gofit.db", "Path to BoltDB")
		static  = cmd.StringOpt("static", "./www", "Path to static assets")
		debug   = cmd.BoolOpt("--debug", true, "Enable Debugging")
	)
	cmd.Action = func() {
		db, err := database.New(*path, *debug)
		FailOnErr(err)
		server.RunServer(db, *pattern, *static)
	}
}

// Load ingests data into the database
func Load(cmd *cli.Cmd) {
	var (
		name     = cmd.StringOpt("n name", "", "Name for this dataset")
		dataType = cmd.StringOpt("t type", "", "Type of data to load")
		dataPath = cmd.StringOpt("p path", "Takeout/", "Path to your raw dataset")
		dbPath   = cmd.StringOpt("d database", "/tmp/gofit.db", "Path to BoltDB")
		debug    = cmd.BoolOpt("debug", true, "Enable Debugging")
	)
	cmd.Action = func() {
		db, err := database.New(*dbPath, *debug)
		FailOnErr(err)
		if *dataType == "" || *name == "" {
			cmd.PrintHelp()
			os.Exit(1)
		}
		switch *dataType {
		case "tcx":
			data, err := tcx.FromDir(*dataPath, *name)
			FailOnErr(err)
			FailOnErr(db.WriteSeries(data.Load()))
		case "csv":
			data, err := csv.FromFile(*dataPath, *name)
			FailOnErr(err)
			FailOnErr(db.WriteSeries(data.Load()))
		default:
			FailOnErr(fmt.Errorf("Unknown datatype %s", *dataType))
		}
	}
}

func Run() {
	app := cli.App("gofit", "GoFit!")
	app.Command("server", "Run the GoFit web UI", Server)
	app.Command("load", "Load New Data", Load)
	app.Run(os.Args)
}
