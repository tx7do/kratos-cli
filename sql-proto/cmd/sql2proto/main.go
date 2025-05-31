package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "sql2proto",
	Short: "SQL to Protobuf code Importer",
	Long:  "sql2proto is a tool to import SQL database schema and generate Protobuf code.",
	Run:   command,
}

var (
	drv           string
	dsn           string
	protoPath     string
	tables        []string
	excludeTables []string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&drv, "drv", "v", "mysql", "Database driver name to use (mysql, postgres, sqlite...)")
	rootCmd.PersistentFlags().StringVarP(&dsn, "dsn", "n", "", `Data source name (connection information), for example:
"mysql://user:pass@tcp(localhost:3306)/dbname"
"postgres://user:pass@host:port/dbname"`)
	rootCmd.PersistentFlags().StringVarP(&protoPath, "proto-path", "s", "./api/protos/", "output path for protobuf schema files")
	rootCmd.PersistentFlags().StringSliceVarP(&tables, "tables", "t", nil, "comma-separated list of tables to inspect (all if empty)")
	rootCmd.PersistentFlags().StringSliceVarP(&excludeTables, "exclude-tables", "e", nil, "comma-separated list of tables to exclude")
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
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute command failed: %v", err)
	}
}
