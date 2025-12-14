package main

import (
	"github.com/spf13/cobra"
	"github.com/tx7do/kratos-cli/gowind/internal/service"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add microservice to an existing project",
	Long:  "Add microservice to the current project. Example: gow add <name>  (shorthand for gow add service <name>)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		switch args[0] {
		case "service", "svc":
			return cmd.Help()
		default:
			service.CmdService.Run(service.CmdService, args)
			return nil
		}
	},
}

func init() {
	addCmd.AddCommand(service.CmdService)
	rootCmd.AddCommand(addCmd)
}
