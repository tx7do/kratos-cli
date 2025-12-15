package main

import (
	"context"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	sqlkratos "github.com/tx7do/kratos-cli/sql-kratos"
)

var rootCmd = &cobra.Command{
	Use:   "sql2kratos",
	Short: "SQL to Kratos microservice code generator",
	Long:  "sql2kratos imports the SQL database schemas and generates Kratos microservice code.",
	Run:   command,
}

var opts sqlkratos.GeneratorOptions

func init() {
	rootCmd.PersistentFlags().StringVarP(&opts.Source, "dsn", "n", "", `Data source name (connection information), for example:
"mysql://user:pass@tcp(localhost:3306)/dbname"
"postgres://user:pass@host:port/dbname"`)

	rootCmd.PersistentFlags().StringVarP(&opts.OutputPath, "output", "o", "./", "output path for the generated code, e.g., './'")

	rootCmd.PersistentFlags().StringVarP(&opts.SourceModuleName, "src-module", "s", "user", "Source module name, for REST service generate, e.g., \"admin\"")
	rootCmd.PersistentFlags().StringVarP(&opts.ModuleName, "module", "m", "admin", "Target module name for the generated code, e.g., 'admin'")
	rootCmd.PersistentFlags().StringVarP(&opts.ModuleVersion, "version", "v", "v1", "Version of the module, e.g., 'v1'")

	rootCmd.PersistentFlags().StringVarP(&opts.ProjectName, "project", "p", "kratos-admin", "Project name for the generated code, e.g., 'kratos-admin'")
	rootCmd.PersistentFlags().StringVarP(&opts.ServiceName, "service", "c", "user", "Service name for the generated code, e.g., 'user'")

	rootCmd.PersistentFlags().StringVarP(&opts.OrmType, "orm", "r", "ent", "ORM type to use (ent, gorm)")

	rootCmd.PersistentFlags().StringSliceVarP(&opts.IncludedTables, "includes", "i", nil, "comma-separated list of tables to inspect (all if empty)")
	rootCmd.PersistentFlags().StringSliceVarP(&opts.ExcludedTables, "excludes", "e", nil, "comma-separated list of tables to exclude")

	rootCmd.PersistentFlags().StringSliceVarP(&opts.Servers, "servers", "g", []string{"grpc"}, "comma-separated list of servers to generate, e.g., \"grpc,rest\"")

	rootCmd.PersistentFlags().BoolVarP(&opts.UseRepo, "repo", "x", true, "use repository pattern")

	rootCmd.PersistentFlags().BoolVarP(&opts.GenerateProto, "gen-proto", "q", true, "enable generate protobuf schema files")
	rootCmd.PersistentFlags().BoolVarP(&opts.GenerateServer, "gen-srv", "w", true, "enable generate server package code")
	rootCmd.PersistentFlags().BoolVarP(&opts.GenerateService, "gen-svc", "a", true, "enable generate service package code")
	rootCmd.PersistentFlags().BoolVarP(&opts.GenerateORM, "gen-orm", "z", true, "enable generate ORM code")
	rootCmd.PersistentFlags().BoolVarP(&opts.GenerateData, "gen-data", "l", true, "enable generate data package code")
	rootCmd.PersistentFlags().BoolVarP(&opts.GenerateMain, "gen-main", "k", true, "enable generate main package code")
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
		log.Println("sql2kratos: dsn must be provided")
		_ = cmd.Help()
		os.Exit(2)
	}

	ctx := context.Background()

	_ = sqlkratos.Generate(ctx, opts)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute command failed: %v", err)
	}
}
