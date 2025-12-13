package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/tx7do/kratos-cli/gowind/internal/project"
	"github.com/tx7do/kratos-cli/gowind/internal/run"
)

var rootCmd = &cobra.Command{
	Use:   "gow",
	Short: "gow CLI",
	Long:  "gow is the CLI for GoWind framework.",
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
