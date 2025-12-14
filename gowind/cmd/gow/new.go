package main

import (
	"github.com/spf13/cobra"

	"github.com/tx7do/kratos-cli/gowind/internal/project"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create new project",
	Long:  "Create new project. Example: gow new <name>  (shorthand for gow new project <name>)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		switch args[0] {
		case "project", "proj":
			return cmd.Help()
		default:
			project.CmdProject.Run(project.CmdProject, args)
			return nil
		}
	},
}

func init() {
	newCmd.AddCommand(project.CmdProject)
	rootCmd.AddCommand(newCmd)
}
