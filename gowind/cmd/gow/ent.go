package main

import (
	"github.com/spf13/cobra"
	"github.com/tx7do/go-wind-toolkit/gowind/internal/ent"
)

var entCmd = &cobra.Command{
	Use:   "ent <service>",
	Short: "manage ent schemas",
	Long:  "Manage ent schemas for services. Examples: gow ent generate admin  or  gow ent add admin User,Group",
	RunE:  ent.RunGenerate,
}

var entGenerateCmd = &cobra.Command{
	Use:   "generate <service>",
	Short: "generate ent code for a service",
	//Args:  cobra.MinimumNArgs(1),
	RunE: ent.RunGenerate,
}

var entAddCmd = &cobra.Command{
	Use:   "add <service> <schemas>",
	Short: "add schema(s) to a service (comma separated, e.g. User,Group)",
	Args:  cobra.MinimumNArgs(2),
	RunE:  ent.RunAdd,
}

func init() {
	//entCmd.AddCommand(entGenerateCmd, entAddCmd)
	rootCmd.AddCommand(entCmd)
}
