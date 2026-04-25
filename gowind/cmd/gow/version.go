package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "v0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gow",
	Long:  `All software has versions. This is gow's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gow version", version)
	},
}
