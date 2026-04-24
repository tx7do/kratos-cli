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

// Run service.
func Run(cmd *cobra.Command, args []string) {
	cmdArgs, _ := pkg.SplitArgs(cmd, args)

	inspector, err := pkg.NewModuleInspectorFromGo("")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return
	}

	// 先在模块根目录运行 `go mod tidy`
	if err = pkg.GoModTidy(cmd.Context(), inspector.Root); err != nil {
		return
	}

	var serviceName string

	if len(cmdArgs) > 0 {
		serviceName = strings.TrimSpace(cmdArgs[0])
		if serviceName == "" {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: service name is required\033[m\n")
			return
		}

		var valid bool
		valid, err = pkg.IsValidServiceName(inspector.Root, serviceName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
			return
		}

		if !valid {
			err = fmt.Errorf("service '%s' does not exist or is not valid (missing cmd/server or configs)", serviceName)
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
			return
		}
	} else {
		// 未指定服务名称，检查当前目录是否为服务目录

		var wd string
		wd, err = os.Getwd()
		if err != nil {
			fmt.Printf("os.Getwd error: %v\n", err)
		}

		var hasCmd, hasConfigs bool
		hasCmd, hasConfigs, err = pkg.HasCmdAndConfigs(wd)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
			return
		}

		//log.Printf("[%s] hasCmd: %v, hasConfigs: %v\n", wd, hasCmd, hasConfigs)

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

	servicePath := path.Join(inspector.Root, "/app/", serviceName, "/service")

	if err = runService(servicePath); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
	}
}

// runService 运行服务，使用命令: go run ./cmd/server -c ./configs。
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

	runArgs := []string{"run", appPathAbs, "-c", configPathAbs}

	if err = g.Run(runArgs...); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return err
	}

	return nil
}
