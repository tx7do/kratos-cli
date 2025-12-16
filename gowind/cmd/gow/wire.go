package main

import (
	"github.com/spf13/cobra"

	"github.com/tx7do/kratos-cli/gowind/internal/wire"
)

var wireCmd = &cobra.Command{
	Use:   "wire",
	Short: "Generate and manage Wire dependency injection code",
	Long: `Run the Wire tool in the specified service directory to generate dependency injection code.
Usage example: wire generate <service>
The command executes in the project's app/<service>/service directory and forwards stdout/stderr to the console.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var wireGenerateCmd = &cobra.Command{
	Use:   "generate <service>",
	Short: "generate wire code for a service",
	//Args:  cobra.MinimumNArgs(1),
	RunE: wire.RunGenerate,
}

func init() {
	wireCmd.AddCommand(wireGenerateCmd)
	rootCmd.AddCommand(wireCmd)
}
