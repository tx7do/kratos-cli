package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/tx7do/kratos-cli/sql-orm"
	"github.com/tx7do/kratos-cli/sql-orm/internal"
)

var rootCmd = &cobra.Command{
	Use:   "sql2orm",
	Short: "SQL to ORM code Importer",
	Long:  "SQL to ORM code Importer is a tool to generate ORM code from SQL database schemas.",
	Run:   command,
}

var opts internal.Options

func init() {
	rootCmd.PersistentFlags().StringVarP(&opts.ORM, "orm", "o", "ent", "ORM type to use (ent, gorm)")
	rootCmd.PersistentFlags().StringVarP(&opts.Driver, "drv", "v", "mysql", "Database driver name to use (mysql, postgres, sqlite...)")
	rootCmd.PersistentFlags().StringVarP(&opts.Source, "dsn", "n", "", `Data source name (connection information), for example:
"mysql://user:pass@tcp(localhost:3306)/dbname"
"postgres://user:pass@host:port/dbname"`)

	rootCmd.PersistentFlags().StringVarP(&opts.SchemaPath, "schema", "s", "./ent/schema/", "output path for schema files (for ent)")
	rootCmd.PersistentFlags().StringVarP(&opts.DaoPath, "dao", "d", "./daos/", "output path for DAO code (for gorm)")
	rootCmd.PersistentFlags().StringSliceVarP(&opts.IncludedTables, "includes", "i", nil, "comma-separated list of tables to inspect (all if empty)")
	rootCmd.PersistentFlags().StringSliceVarP(&opts.ExcludedTables, "excludes", "e", nil, "comma-separated list of tables to exclude")
}

func parseDSN(url string) (string, string, error) {
	a := strings.SplitN(url, "://", 2)
	if len(a) != 2 {
		return "", "", fmt.Errorf(`failed to parse dsn: "%s"`, url)
	}
	return a[0], a[1], nil
}

// countFlags 统计显式设置的标志数量
func countFlags(cmd *cobra.Command) int {
	count := 0
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Changed {
			count++
		}
	})
	return count
}

func command(cmd *cobra.Command, _ []string) {
	if countFlags(cmd) == 0 {
		_ = cmd.Help()
		return
	}

	if opts.Source == "" {
		log.Println("sql2orm: dsn must be provided")
		_ = cmd.Help()
		os.Exit(2)
	}
	if opts.ORM == "" {
		opts.ORM = "ent"
	}

	ctx := context.Background()

	_ = sqlorm.Importer(
		ctx,
		opts.ORM,
		&opts.Driver, &opts.Source,
		&opts.SchemaPath, &opts.DaoPath,
		opts.IncludedTables, opts.ExcludedTables,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute command failed: %v", err)
	}
}
