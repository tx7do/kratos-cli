package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tx7do/kratos-cli/sql-importer/internal/ent/entimport"
	"log"
	"os"
	"strings"
)

var (
	tablesFlag        tables
	excludeTablesFlag tables
)

func init() {
	flag.Var(&tablesFlag, "tables", "comma-separated list of tables to inspect (all if empty)")
	flag.Var(&excludeTablesFlag, "exclude-tables", "comma-separated list of tables to exclude")
}

func main() {
	dsn := flag.String("dsn", "",
		`data source name (connection information), for example:
"mysql://user:pass@tcp(localhost:3306)/dbname"
"postgres://user:pass@host:port/dbname"`)
	schemaPath := flag.String("schema-path", "./ent/schema", "output path for ent schema")
	flag.Parse()

	if *dsn == "" {
		log.Println("entimport: data source name (dsn) must be provided")
		flag.Usage()
		os.Exit(2)
	}

	ctx := context.Background()

	_ = entimport.Importer(ctx, dsn, schemaPath, tablesFlag, excludeTablesFlag)
}

type tables []string

func (t *tables) String() string {
	return fmt.Sprint(*t)
}

func (t *tables) Set(s string) error {
	*t = strings.Split(s, ",")
	return nil
}
