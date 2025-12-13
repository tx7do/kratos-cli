package main

import (
	"github.com/spf13/cobra"

	"github.com/tx7do/kratos-cli/gowind/internal/project"
	"github.com/tx7do/kratos-cli/gowind/internal/service"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create new project or service",
	Long:  "Create new project or service. Example: gow new <name>  (shorthand for gow new project <name>)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		// 当第一个参数是子命令名时，交给 cobra 正常处理（返回 Help，让 cobra 显示用法）
		switch args[0] {
		case "project", "proj", "service", "svc":
			return cmd.Help()
		default:
			// 把 args 转发给 project 子命令（等同于 gow new project <name>）
			project.CmdProject.Run(project.CmdProject, args)
			return nil
		}
	},
}

func init() {
	// 将子命令注册到 newCmd
	newCmd.AddCommand(project.CmdProject, service.CmdService)
	// 然后把 newCmd 注册到根命令（假设有 rootCmd）
	rootCmd.AddCommand(newCmd)
}
