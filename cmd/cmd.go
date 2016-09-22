package cmd

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/kevinschoon/fit/database"
	"github.com/kevinschoon/fit/loaders"
	"github.com/kevinschoon/fit/loaders/csv"
	"github.com/kevinschoon/fit/server"
	"os"
)

const FitVersion string = "0.0.1"

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
		demo    = cmd.BoolOpt("demo", false, "Run in Demo Mode")
		debug   = cmd.BoolOpt("--debug", true, "Enable Debugging")
	)
	cmd.Action = func() {
		db, err := database.New(*path, *debug)
		FailOnErr(err)
		server.RunServer(db, *pattern, *static, FitVersion, *demo)
	}
}

// Load ingests data into the database
func Load(cmd *cli.Cmd) {
	var (
		path     = cmd.StringArg("PATH", "", "Path to your raw dataset")
		name     = cmd.StringOpt("n name", "", "Name for this dataset")
		ft       = cmd.StringOpt("t type", "", "Type of data to load")
		dtIndex  = cmd.IntOpt("dtIndex", 0, "Column to extract time.Time from")
		dtFormat = cmd.StringOpt("dtFormat", "", "Format to extract time.Time with")
		//server   = cmd.StringOpt("server", "http://localhost:8000", "Fit API server")
		series = cmd.BoolOpt("series", false, "Dump series information to stdout")
		values = cmd.BoolOpt("values", false, "Dump values information to stdout")
	)
	cmd.Action = func() {
		opts := &loaders.Options{
			Name: *name,
			Path: *path,
			Type: loaders.FileTypeByName(*ft),
			CSVOptions: &csv.Options{
				DTIndex:  *dtIndex,
				DTFormat: *dtFormat,
			},
		}
		loader, err := loaders.Load(opts)
		FailOnErr(err)
		defer loader.Close()
		switch {
		case *series:
			FailOnErr(loaders.Stdout(opts.Name, false, loader))
		case *values:
			FailOnErr(loaders.Stdout(opts.Name, true, loader))
		}
	}
}

func Run() {
	app := cli.App("fit", "Fit is a toolkit for exploring numerical data")
	app.Command("server", "Run the Fit web server", Server)
	app.Command("load", "Load and transform new data", Load)
	app.Version("v version", FitVersion)
	app.Run(os.Args)
}
