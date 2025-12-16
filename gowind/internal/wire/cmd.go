package wire

import (
	"github.com/spf13/cobra"
	"github.com/tx7do/kratos-cli/gowind/internal/pkg"
)

// CmdRun run project command.
var CmdRun = &cobra.Command{
	Use:   "wire",
	Short: "generate wire code",
	Long:  "Generate wire code. Example: gowind wire admin",
	Run:   Run,
}

var serviceName string

// Run service.
func Run(cmd *cobra.Command, args []string) {
	cmdArgs, _ := pkg.SplitArgs(cmd, args)
}
