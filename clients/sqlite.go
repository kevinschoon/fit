package clients

import (
	"bytes"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/kevinschoon/fit/types"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"text/template"
)

const (
	createStmt = "CREATE TABLE IF NOT EXISTS {{.Name}} (id integer not null primary key,{{range .Columns}}{{.}}float,{{end}})"
)

// SQLClient implements a client interface
// for reading and writing datasets
type SQLClient struct {
	db *sqlx.DB
}

func (c *SQLClient) exists(name string) bool {
	return c.db.Get("SELECT name FROM sqlite_master WHERE type='table' AND name='?'", name) == nil
}

func (c *SQLClient) Datasets() (datasets []*types.Dataset, err error) {
	rows, err := c.db.Queryx("SELECT * FROM datasets")
	if err != nil {
		return nil, err
	}
	datasets = make([]*types.Dataset, 0)
	for rows.Next() {
		dataset := &types.Dataset{}
		if err := rows.StructScan(dataset); err != nil {
			return nil, err
		}
		datasets = append(datasets, dataset)
	}
	return datasets, nil
}

func (c *SQLClient) Write(ds *types.Dataset) (err error) {
	stmt := "CREATE TABLE IF NOT EXISTS {{.Name}} 
	(id integer not null primary key,{{range .Columns}}{{.}}float,{{end}})"
	tmpl := template.Must(template.New("").Parse(stmt))
	buf := bytes.Buffer{}
	if err := tmpl.Execute(buf, ds); err != nil {
		return err
	}
	if _, err = c.db.Exec(buf.String()); err != nil {
		return err
	}
	_, err = c.db.Exec("INSERT INTO datasets () VALUES (?, ?, ?, ?)",
		ds.Name, strings.Join(ds.Columns, ","), ds.Stats.Rows, ds.Stats.Columns)
	if err != nil {
		return err
	}
	return nil
}

func (c *SQLClient) Delete(name string) (err error) {
	_, err = c.db.Exec("DROP TABLE ?", name)
	return err
}

func (c *SQLClient) Query(query *types.SQLQuery) (ds *types.Dataset, err error) {
	return nil, nil
}

// NewSQLClient returns a new SQLite client
// TODO: Add close() method to client interface
func NewSQLClient(path string) (*SQLClient, error) {
	db, err := sqlx.Connect("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return &SQLClient{
		db: db,
	}, nil
}
