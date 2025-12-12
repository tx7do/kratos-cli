package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/tx7do/kratos-cli/gowind/internal/project"
)

var rootCmd = &cobra.Command{
	Use:   "gowind",
	Short: "gowind: Go microservices toolkit",
	Long:  `gowind: Go microservices toolkit.`,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
