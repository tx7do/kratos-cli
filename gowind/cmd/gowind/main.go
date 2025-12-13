package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/tx7do/kratos-cli/gowind/internal/project"
	"github.com/tx7do/kratos-cli/gowind/internal/run"
)

var rootCmd = &cobra.Command{
	Use:   "gowind",
	Short: "gowind: GoWind microservices toolkit",
	Long:  `gowind: GoWind microservices toolkit.`,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(run.CmdRun)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
