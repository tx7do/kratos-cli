package main

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	sqlproto "github.com/tx7do/kratos-cli/sql-proto"
	"github.com/tx7do/kratos-cli/sql-proto/internal"
)

var rootCmd = &cobra.Command{
	Use:   "sql2proto",
	Short: "SQL to Protobuf code generator",
	Long:  "sql2proto imports the SQL database schemas and generates Protobuf code for use in Kratos microservices.",
	Run:   command,
}

var opts internal.Options

func init() {
	//rootCmd.PersistentFlags().StringVarP(&opts.Driver, "drv", "v", "mysql", "Database driver name to use (mysql, postgres, sqlite...)")
	rootCmd.PersistentFlags().StringVarP(&opts.Source, "dsn", "n", "", `Data source name (connection information), for example:
"mysql://user:pass@tcp(localhost:3306)/dbname"
"postgres://user:pass@host:port/dbname"`)
	rootCmd.PersistentFlags().StringVarP(&opts.OutputPath, "output", "o", "./api/protos/", "output path for protobuf schema files")
	rootCmd.PersistentFlags().StringVarP(&opts.SourceModule, "src-module", "s", "user", "Source module name, for REST service generate, e.g., \"admin\"")
	rootCmd.PersistentFlags().StringVarP(&opts.Module, "module", "m", "admin", "module name for the generated code, e.g., 'admin'")
	rootCmd.PersistentFlags().StringVarP(&opts.Version, "version", "v", "v1", "Version of the module, e.g., 'v1'")
	rootCmd.PersistentFlags().StringVarP(&opts.Service, "type", "t", "grpc", "generate RPC service type, \"rest\" for REST service, \"grpc\" for gRPC service")
	rootCmd.PersistentFlags().StringSliceVarP(&opts.IncludedTables, "includes", "i", nil, "comma-separated list of tables to inspect (all if empty)")
	rootCmd.PersistentFlags().StringSliceVarP(&opts.ExcludedTables, "excludes", "e", nil, "comma-separated list of tables to exclude")
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
		log.Println("sql2proto: dsn must be provided")
		_ = cmd.Help()
		os.Exit(2)
	}

	ctx := context.Background()

	_ = sqlproto.Convert(
		ctx,
		&opts.Source, &opts.OutputPath,
		&opts.Module,
		&opts.SourceModule,
		&opts.Version,
		&opts.Service,
		opts.IncludedTables, opts.ExcludedTables,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute command failed: %v", err)
	}
}
