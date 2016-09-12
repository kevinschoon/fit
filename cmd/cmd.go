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

func Server(cmd *cli.Cmd) {
	pattern := *cmd.StringOpt("pattern", ":8000", "IP and port pattern to listen on")
	dbPath := *cmd.StringOpt("p path", "/tmp/gofit.db", "Path to BoltDB")
	static := *cmd.StringOpt("static", "./www", "Path to static assets")
	cmd.Action = func() {
		db, err := database.New(dbPath)
		FailOnErr(err)
		server.RunServer(db, pattern, static)
	}
}

func Load(cmd *cli.Cmd) {
	name := cmd.StringOpt("n name", "", "Name for this dataset")
	dataType := cmd.StringOpt("t type", "", "Type of data to load")
	dataPath := cmd.StringOpt("p path", "Takeout/", "Path to your raw dataset")
	dbPath := cmd.StringOpt("d database", "/tmp/gofit.db", "Path to BoltDB")
	cmd.Action = func() {
		db, err := database.New(*dbPath)
		FailOnErr(err)
		fmt.Println(*name, *dataType, *dataPath, *dbPath)
		switch *dataType {
		case "tcx":
			data, err := tcx.FromDir(*dataPath)
			FailOnErr(err)
			FailOnErr(db.Write(*name, data.Load()))
		case "csv":
			data, err := csv.FromFile(*dataPath)
			FailOnErr(err)
			FailOnErr(db.Write(*name, data.Load()))
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
