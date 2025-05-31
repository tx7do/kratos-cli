package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/tx7do/kratos-cli/config-exporter"
	"github.com/tx7do/kratos-cli/config-exporter/internal"
)

var rootCmd = &cobra.Command{
	Use:   "cfgexp",
	Short: "Config Exporter",
	Long:  "Config Exporter is a tool to export configuration from remote services like Consul or Etcd to local files.",
	Run:   command,
}

var opts internal.Options

func init() {
	rootCmd.PersistentFlags().StringVarP((*string)(&opts.Service), "type", "t", "consul", "remote config service name (consul, etcd, etc.)")
	rootCmd.PersistentFlags().StringVarP(&(opts.Endpoint), "addr", "a", "127.0.0.1:8500", "remote config service address")
	rootCmd.PersistentFlags().StringVarP(&(opts.ProjectName), "proj", "p", "", "project name, this name is used to key prefix in remote config service")
	rootCmd.PersistentFlags().StringVarP(&(opts.ProjectRoot), "root", "r", "./", "project root dir")
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

	_ = cfgexp.Export(string(opts.Service), opts.Endpoint, opts.ProjectName, opts.ProjectRoot)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute command failed: %v", err)
	}
}
