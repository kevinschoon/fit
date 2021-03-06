package cmd

import (
	"encoding/json"
	"fmt"
	mtx "github.com/gonum/matrix/mat64"
	"github.com/gosuri/uitable"
	"github.com/jawher/mow.cli"

	"github.com/kevinschoon/fit/clients"
	"github.com/kevinschoon/fit/loader"
	"github.com/kevinschoon/fit/parser"
	"github.com/kevinschoon/fit/server"
	"github.com/kevinschoon/fit/types"
	"os"
)

const FitVersion string = "0.0.1"

var (
	app = cli.App("fit", "Fit is a toolkit for exploring, extracting, and transforming datasets")

	dbPath  = app.StringOpt("d db", "", "Path to a BoltDB database, default: /tmp/fit.db")
	apiURL  = app.StringOpt("s server", "", "Fit API server, default: http://127.0.0.1:8000")
	asHuman = app.BoolOpt("h human", true, "output data as human readable text")
	asJSON  = app.BoolOpt("j json", false, "output data in JSON format")
)

func FailOnErr(err error) {
	if err != nil {
		fmt.Println("ERROR:", err.Error())
		os.Exit(1)
	}
}

func GetClient(pref string) types.Client {
	switch {
	case *dbPath != "":
		client, err := clients.NewBoltClient(*dbPath)
		FailOnErr(err)
		return client
	case *apiURL != "":
		client, err := clients.NewHTTPClient(*apiURL)
		FailOnErr(err)
		return client
	}
	if pref == "db" {
		client, err := clients.NewBoltClient("/tmp/fit.db")
		FailOnErr(err)
		return client
	}
	client, err := clients.NewHTTPClient("http://127.0.0.1:8000")
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
			server.RunServer(GetClient("db"), *pattern, *static, FitVersion, *demo)
		}
	})

	app.Command("load", "load a dataset into BoltDB", func(cmd *cli.Cmd) {
		cmd.Spec = "[[-n] [-s]][-p...][-c...] PATH"
		var (
			name       = cmd.StringOpt("n name", "", "name of this dataset")
			path       = cmd.StringArg("PATH", "", "File path")
			parserArgs = cmd.StringsOpt("p parser", []string{}, "parsers to apply")
			sheet      = cmd.StringOpt("s sheet", "", "name of the sheet to load with XLS file")
			columns    = cmd.StringsOpt("c column", []string{}, "column names")
		)
		cmd.Action = func() {
			parsers, err := parser.ParsersFromArgs(*parserArgs)
			FailOnErr(err)
			opts := loader.Options{
				Name:    *name,
				Path:    *path,
				Parsers: parsers,
				Columns: *columns,
				Sheet:   *sheet,
			}
			ds, err := loader.ReadPath(opts)
			FailOnErr(err)
			FailOnErr(GetClient("").Write(ds))
		}
	})

	app.Command("ls", "list datasets loaded into the database with their columns", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			datasets, err := GetClient("").Datasets()
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
			FailOnErr(GetClient("").Delete(*name))
		}
	})

	app.Command("query", "Query values from one or more datasets", func(cmd *cli.Cmd) {
		var (
			queryArgs = cmd.StringsArg("QUERY", []string{}, "Query parameters")
			lines     = cmd.IntOpt("n lines", 10, "number of rows to output")
			grouping  = cmd.StringOpt("g grouping", "", "grouping to apply to the resulting matrix")
			function  = cmd.StringOpt("f function", "avg", "function to apply when grouping")
		)
		cmd.LongDesc = `Query values from one or more stored datasets. Values from different 
datasets can be joined together by specifying multiple query parameters.

Example:

fit query -n 10 -g Duration,0,1m -f avg "Dataset1,fuu" "Dataset2,bar,baz"
`
		cmd.Spec = "[OPTIONS] QUERY..."
		cmd.Action = func() {
			if len(*queryArgs) == 0 {
				cmd.PrintLongHelp()
				os.Exit(1)
			}
			ds, err := GetClient("").Query(types.NewQuery(*queryArgs, *function, *grouping))
			FailOnErr(err)
			if ds.Len() > 0 {
				switch {
				case *asJSON:
					ds.WithValues = true
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
