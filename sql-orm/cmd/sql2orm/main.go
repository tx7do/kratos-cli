package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/tx7do/kratos-cli/sql-orm/internal/ent/entimport"
	"github.com/tx7do/kratos-cli/sql-orm/internal/gorm"
)

var rootCmd = &cobra.Command{
	Use:   "sql2orm",
	Short: "SQL to ORM code Importer",
	Long:  "SQL to ORM code Importer is a tool to generate ORM code from SQL database schemas.",
	Run:   command,
}

var (
	orm           string
	drv           string
	dsn           string
	schemaPath    string
	daoPath       string
	tables        []string
	excludeTables []string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&orm, "orm", "o", "ent", "ORM type to use (ent, gorm)")
	rootCmd.PersistentFlags().StringVarP(&drv, "drv", "v", "mysql", "Database driver name to use (mysql, postgres, sqlite...)")
	rootCmd.PersistentFlags().StringVarP(&dsn, "dsn", "n", "", `Data source name (connection information), for example:
"mysql://user:pass@tcp(localhost:3306)/dbname"
"postgres://user:pass@host:port/dbname"`)

	rootCmd.PersistentFlags().StringVarP(&schemaPath, "schema-path", "s", "./ent/schema/", "output path for schema")
	rootCmd.PersistentFlags().StringVarP(&daoPath, "dao-path", "d", "./daos/", "output path for DAO code (for gorm)")
	rootCmd.PersistentFlags().StringSliceVarP(&tables, "tables", "t", nil, "comma-separated list of tables to inspect (all if empty)")
	rootCmd.PersistentFlags().StringSliceVarP(&excludeTables, "exclude-tables", "e", nil, "comma-separated list of tables to exclude")
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

	if dsn == "" {
		log.Println("sql2orm: dsn must be provided")
		_ = cmd.Help()
		os.Exit(2)
	}
	if orm == "" {
		orm = "ent"
	}

	ctx := context.Background()

	switch strings.ToLower(strings.TrimSpace(orm)) {
	case "ent":
		_ = entimport.Importer(ctx, &dsn, &schemaPath, tables, excludeTables)

	case "gorm":
		_ = gorm.Importer(ctx, &drv, &dsn, &schemaPath, &daoPath, tables, excludeTables)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute command failed: %v", err)
	}
}
