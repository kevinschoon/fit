package cmd

import (
	"fmt"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/jawher/mow.cli"
	"github.com/kevinschoon/fit/loader"
	"github.com/kevinschoon/fit/server"
	"github.com/kevinschoon/fit/store"
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
		pattern = cmd.StringOpt("pattern", "127.0.0.1:8000", "IP and port pattern to listen on")
		path    = cmd.StringOpt("p path", "/tmp/fit.db", "Path to BoltDB")
		static  = cmd.StringOpt("static", "./www", "Path to static assets")
		demo    = cmd.BoolOpt("demo", false, "Run in Demo Mode")
	)
	cmd.Action = func() {
		db, err := store.NewDB(*path)
		FailOnErr(err)
		server.RunServer(db, *pattern, *static, FitVersion, *demo)
	}
}

func Ls(cmd *cli.Cmd) {
	var dbPath = cmd.StringOpt("d db", "/tmp/fit.db", "Write to BoltDB path")
	cmd.Action = func() {
		db, err := store.NewDB(*dbPath)
		FailOnErr(err)
		datasets, err := db.Datasets()
		FailOnErr(err)
		for _, dataset := range datasets {
			fmt.Printf("%s-%s\n", dataset.Name, dataset.Columns)
		}
	}
}

func Load(cmd *cli.Cmd) {
	var (
		name       = cmd.StringOpt("n name", "", "Name for this dataset")
		parserOpts = cmd.StringsOpt("p parser", []string{}, "Parsers to apply")
		path       = cmd.StringArg("PATH", "", "Path to your CSV dataset")
		dbPath     = cmd.StringOpt("d db", "/tmp/fit.db", "Write to BoltDB path")
	)
	cmd.Action = func() {
		if *name == "" {
			FailOnErr(fmt.Errorf("You must specify a name"))
		}
		db, err := store.NewDB(*dbPath)
		FailOnErr(err)
		parsers, err := loader.ParsersFromArgs(*parserOpts)
		FailOnErr(err)
		reader, err := loader.NewCSV(&loader.CSVOptions{Path: *path, Parsers: parsers})
		FailOnErr(err)
		m, err := store.ReadMatrix(reader)
		FailOnErr(err)
		FailOnErr(db.Write(&store.Dataset{Name: *name, Columns: reader.Columns()}, m))
	}
}

func Cat(cmd *cli.Cmd) {
	var (
		parserOpts = cmd.StringsOpt("p parser", []string{}, "Parsers to apply")
		path       = cmd.StringArg("PATH", "", "Path to your CSV dataset")
	)
	cmd.Action = func() {
		parsers, err := loader.ParsersFromArgs(*parserOpts)
		FailOnErr(err)
		reader, err := loader.NewCSV(&loader.CSVOptions{Path: *path, Parsers: parsers})
		FailOnErr(err)
		m, err := store.ReadMatrix(reader)
		FailOnErr(err)
		fmt.Printf("\n%s\n", reader.Columns())
		fmt.Printf("\n%v\n\n", mtx.Formatted(m, mtx.Prefix("  "), mtx.Excerpt(5)))
	}
}

func Run() {
	app := cli.App("fit", "Fit is a toolkit for exploring numerical data")
	app.Command("server", "Run the Fit web server", Server)
	app.Command("load", "load a dataset into BoltDB", Load)
	app.Command("cat", "write a dataset to stdout or perform transformations on it", Cat)
	app.Command("ls", "list datasets loaded into the database with their columns", Ls)
	app.Version("v version", FitVersion)
	app.Run(os.Args)
}
