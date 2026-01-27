package main

import (
	"github.com/spf13/cobra"

	"github.com/tx7do/kratos-cli/gowind/internal/buf"
)

var bufCmd = &cobra.Command{
	Use:   "api",
	Short: "manage proto and buf files",
	Long:  "Manage proto and buf files for services.",
	Args:  cobra.NoArgs,
	RunE:  buf.RunGenerate,
}

func init() {
	rootCmd.AddCommand(bufCmd)
}
