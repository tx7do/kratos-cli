package run

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tx7do/kratos-cli/gowind/internal/pkg"
)

// CmdRun run project command.
var CmdRun = &cobra.Command{
	Use:   "run",
	Short: "Run service project",
	Long:  "Run service project. Example: gowind run admin",
	Run:   Run,
}

var serviceName string

// Run service.
func Run(cmd *cobra.Command, args []string) {
	cmdArgs, _ := pkg.SplitArgs(cmd, args)

	if len(cmdArgs) > 0 {
		serviceName = strings.TrimSpace(cmdArgs[0])
		if serviceName == "" {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: service name is required\033[m\n")
			return
		}
	} else {
		// 未指定服务名称，检查当前目录是否为服务目录

		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("os.Getwd error: %v\n", err)
		}

		hasCmd, hasConfigs, err := pkg.HasCmdAndConfigs(wd)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
			return
		}

		//fmt.Printf("[%s] hasCmd: %v, hasConfigs: %v\n", wd, hasCmd, hasConfigs)

		if hasCmd && hasConfigs {
			// 当前目录即为服务目录
			if err = runService(wd); err != nil {
				return
			}
			return
		}

		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: this is not a valid service folder\033[m\n")
		return
	}

	inspector, err := pkg.NewModuleInspectorFromGo("")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return
	}

	servicePath := path.Join(inspector.Root, "/app/", serviceName, "/service")

	if err = runService(servicePath); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
	}
}

// runService 运行服务，使用命令: go run ./cmd/server -conf ./configs。
func runService(serviceWorkPath string) error {
	// 使用 pkg.NewGoCmd 执行 go run . [programArgs...]
	g := pkg.NewGoCmd(serviceWorkPath)
	g.Stdout = os.Stdout
	g.Stderr = os.Stderr

	// 构建并规范化路径
	appPath := filepath.Join(serviceWorkPath, "cmd", "server")
	appPathAbs, err := filepath.Abs(appPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return err
	}
	appPathAbs = filepath.Clean(appPathAbs)

	configPath := filepath.Join(serviceWorkPath, "configs")
	configPathAbs, err := filepath.Abs(configPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return err
	}
	configPathAbs = filepath.Clean(configPathAbs)

	runArgs := []string{"run", appPathAbs, "-conf", configPathAbs}

	if err = g.Run(runArgs...); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return err
	}

	return nil
}
