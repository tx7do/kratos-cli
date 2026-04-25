package main

import (
	"github.com/spf13/cobra"

	"github.com/tx7do/go-wind-toolkit/gowind/internal/wire"
)

var wireCmd = &cobra.Command{
	Use:   "wire <service>",
	Short: "Generate Wire dependency injection code",
	Long: `Run the Wire tool in the specified service directory to generate dependency injection code.
Usage example: wire <service>
The command executes in the project's app/<service>/service directory and forwards stdout/stderr to the console.`,
	RunE: wire.RunGenerate,
}

func init() {
	rootCmd.AddCommand(wireCmd)
}
