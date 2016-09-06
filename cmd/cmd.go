package cmd

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/kevinschoon/gofit/database"
	"github.com/kevinschoon/gofit/server"
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
	pattern := *cmd.StringOpt("pattern", ":8000", "IP and port pattern to listen on")
	dbPath := *cmd.StringOpt("p path", "/tmp/gofit.db", "Path to SQLite DB")
	cmd.Action = func() {
		_, err := database.New(dbPath, &tcx.Activity{}, tcx.Lap{}, tcx.Track{}, tcx.Trackpoint{})
		FailOnErr(err)
		server.RunServer(pattern, dbPath)
	}
}

func Load(cmd *cli.Cmd) {
	tcxPath := *cmd.StringOpt("t tcxPath", "Takeout/", "Path to your TCX Data")
	dbPath := *cmd.StringOpt("p path", "/tmp/gofit.db", "Path to SQLite DB")
	cmd.Action = func() {
		db, err := database.New(dbPath, &tcx.Activity{}, tcx.Lap{}, tcx.Track{}, tcx.Trackpoint{})
		FailOnErr(err)
		loader := database.TCXLoader{Path: tcxPath}
		FailOnErr(loader.Load())
		FailOnErr(database.Write(db, loader))
	}
}

func Run() {
	app := cli.App("gofit", "GoFit!")
	app.Command("server", "Run the GoFit web UI", Server)
	app.Command("load", "Load TCX Data", Load)
	app.Run(os.Args)
}
