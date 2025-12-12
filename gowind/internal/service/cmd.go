package service

import "github.com/spf13/cobra"

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "service",
	Short: "Create a project template",
	Long:  "Create a project using the repository template. Example: gowind service user",
	Run:   run,
}

func run(_ *cobra.Command, args []string) {

}
