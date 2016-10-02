package cmd

import (
	"encoding/json"
	"fmt"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	"github.com/kevinschoon/fit/loader"
	"github.com/kevinschoon/fit/parser"
	"github.com/kevinschoon/fit/server"
	"github.com/kevinschoon/fit/store"
	"os"
)

const FitVersion string = "0.0.1"

var (
	app = cli.App("fit", "Fit is a toolkit for exploring, extracting, and transforming datasets")

	dbPath  = app.StringOpt("d db", "/tmp/fit.db", "Path to a BoltDB database")
	asHuman = app.BoolOpt("h human", true, "output data as human readable text")
	asJSON  = app.BoolOpt("j json", false, "output data in JSON format")
)

func FailOnErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}

func GetDB() *store.DB {
	db, err := store.NewDB(*dbPath)
	FailOnErr(err)
	return db
}

func Run() {
	app.Version("v version", FitVersion)

	app.Command("server", "Run the Fit web server", func(cmd *cli.Cmd) {
		var (
			pattern = cmd.StringOpt("pattern", "127.0.0.1:8000", "IP and port pattern to listen on")
			static  = cmd.StringOpt("static", "./www", "Path to static assets")
			demo    = cmd.BoolOpt("demo", false, "Run in Demo Mode")
		)
		cmd.Action = func() {
			db := GetDB()
			server.RunServer(db, *pattern, *static, FitVersion, *demo)
		}
	})

	app.Command("load", "load a dataset into BoltDB", func(cmd *cli.Cmd) {
		cmd.Spec = "[-n][-p...] PATH"
		var (
			name    = cmd.StringOpt("n name", "", "name of this dataset")
			path    = cmd.StringArg("PATH", "", "File path")
			parsers = cmd.StringsOpt("p parser", []string{}, "parsers to apply")
		)
		cmd.Action = func() {
			p, err := parser.ParsersFromArgs(*parsers)
			FailOnErr(err)
			ds, err := loader.ReadPath(*name, *path, loader.NONE, p)
			FailOnErr(err)
			FailOnErr(GetDB().Write(ds))
		}
	})

	app.Command("ls", "list datasets loaded into the database with their columns", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			db := GetDB()
			defer db.Close()
			datasets, err := db.Datasets()
			FailOnErr(err)
			switch {
			case *asJSON:
				raw, err := json.Marshal(datasets)
				FailOnErr(err)
				fmt.Println(string(raw))
			default:
				tbl := uitable.New()
				tbl.AddRow("NAME", "ROWS", "COLS", "COLUMNS")
				for _, dataset := range datasets {
					tbl.AddRow(dataset.Name, fmt.Sprintf("%d", dataset.Stats.Rows), fmt.Sprintf("%d", dataset.Stats.Columns), dataset.Columns)
				}
				fmt.Println(tbl)
			}
		}
	})

	app.Command("show", "Show values from one or more stored datasets", func(cmd *cli.Cmd) {
		var (
			queryArgs = cmd.StringsOpt("q query", []string{}, "Query parameters")
			lines     = cmd.IntOpt("n lines", 10, "number of rows to output")
		)
		cmd.LongDesc = `Show values from one or more stored datasets. Values from different 
datasets can be joined together by specifying multiple query parameters.

Example:

fit show -q Dataset1,fuu -q Dataset2,bar,baz -n 5
`
		cmd.Spec = "[-q...][-n]"
		cmd.Action = func() {
			if len(*queryArgs) == 0 {
				cmd.PrintLongHelp()
				os.Exit(1)
			}
			db := GetDB()
			ds, err := db.Query(store.NewQueries(*queryArgs))
			FailOnErr(err)
			if ds.Len() > 0 {
				switch {
				case *asJSON:
					raw, err := json.Marshal(ds)
					FailOnErr(err)
					fmt.Println(string(raw))
				default:
					fmt.Printf("\n%s\n", ds.Columns)
					if *lines > 0 {
						fmt.Printf("\n%v\n\n", mtx.Formatted(ds.Mtx, mtx.Prefix("  "), mtx.Excerpt(*lines)))
					} else {
						fmt.Printf("\n%v\n\n", mtx.Formatted(ds.Mtx, mtx.Prefix("  ")))
					}
				}
			}
		}
	})

	FailOnErr(app.Run(os.Args))
}
