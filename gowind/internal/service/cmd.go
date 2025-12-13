package service

import (
	"github.com/spf13/cobra"
)

// CmdService represents the service command
var CmdService = &cobra.Command{
	Use:     "service [name]",
	Aliases: []string{"svc"},
	Short:   "create a new service scaffold",
	Long:    "Create a new microservice inside the current workspace. Example: gow new service usersvc",
	Args:    cobra.ExactArgs(1),
	Run:     run,
}

func run(_ *cobra.Command, args []string) {

}
