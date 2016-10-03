package cmd

import (
	"encoding/json"
	"fmt"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"
	"github.com/kevinschoon/fit/client"
	"github.com/kevinschoon/fit/loader"
	"github.com/kevinschoon/fit/parser"
	"github.com/kevinschoon/fit/server"
	"github.com/kevinschoon/fit/store"
	"os"
)

const FitVersion string = "0.0.1"

var (
	app = cli.App("fit", "Fit is a toolkit for exploring, extracting, and transforming datasets")

	dbPath  = app.StringOpt("d db", "", "Path to a BoltDB database")
	apiURL  = app.StringOpt("s server", "http://127.0.0.1:8000", "Fit API server")
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

func GetClient() *client.Client {
	client, err := client.NewClient(*apiURL)
	FailOnErr(err)
	return client
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
			switch {
			case *dbPath != "":
				FailOnErr(GetDB().Write(ds))
			default:
				FailOnErr(GetClient().Write(ds))
			}
		}
	})

	app.Command("ls", "list datasets loaded into the database with their columns", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			var (
				datasets []*store.Dataset
				err      error
			)
			switch {
			case *dbPath != "":
				datasets, err = GetDB().Datasets()
			default:
				datasets, err = GetClient().Datasets()
			}
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

	app.Command("rm", "Delete a dataset", func(cmd *cli.Cmd) {
		var name = cmd.StringArg("NAME", "", "Name of the dataset to delete")
		cmd.Action = func() {
			if *name == "" {
				cmd.PrintLongHelp()
				os.Exit(1)
			}
			switch {
			case *dbPath != "":
				FailOnErr(GetDB().Delete(*name))
			default:
				FailOnErr(GetClient().Delete(*name))
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
			var (
				ds  *store.Dataset
				err error
			)
			switch {
			case *dbPath != "":
				ds, err = GetDB().Query(store.NewQueries(*queryArgs))
			default:
				ds, err = GetClient().Query(store.NewQueries(*queryArgs))
			}
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
