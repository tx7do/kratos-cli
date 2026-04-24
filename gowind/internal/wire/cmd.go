package wire

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tx7do/kratos-cli/gowind/internal/pkg"
)

// RunGenerate 用于 cobra 的 RunE：校验参数并在服务目录执行 wire 命令。
func RunGenerate(cmd *cobra.Command, args []string) error {
	cmdArgs, _ := pkg.SplitArgs(cmd, args)

	inspector, err := pkg.NewModuleInspectorFromGo("")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return err
	}

	// 先在模块根目录运行 `go mod tidy`
	if err = pkg.GoModTidy(cmd.Context(), inspector.Root); err != nil {
		return err
	}

	var serviceName string

	if len(cmdArgs) > 0 {
		serviceName = strings.TrimSpace(cmdArgs[0])
		if serviceName == "" {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: service name is required\033[m\n")
			return fmt.Errorf("service name is required")
		}

		var valid bool
		valid, err = pkg.IsValidServiceName(inspector.Root, serviceName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
			return err
		}

		if !valid {
			err = fmt.Errorf("service '%s' does not exist or is not valid (missing cmd/server or configs)", serviceName)
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
			return err
		}
	} else {
		// 未指定服务名称，检查当前目录是否为服务目录

		var wd string
		wd, err = os.Getwd()
		if err != nil {
			fmt.Printf("os.Getwd error: %v\n", err)
		}

		serviceName, err = pkg.ExtractServiceName(inspector.Root, wd)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
			return err
		}
	}

	// 如果没有传入 service 或者为空，则遍历所有 service
	if len(serviceName) == 0 {
		return generateWireAllService(cmd.Context(), inspector.Root)
	}

	servicePath := filepath.Join(inspector.Root, "app", serviceName, "service", "cmd", "server")
	err = generateWire(cmd.Context(), servicePath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: generate for service %s failed: %v\033[m\n", serviceName, err)
		return err
	}

	_, _ = fmt.Fprintf(os.Stdout, "\033[32mSUCCESS: wire generated for service '%s'\033[m\n", serviceName)

	return nil
}

// generateWireAllService 在模块根目录下查找 app 目录，并对其中每个包含 service/cmd/server 的服务目录执行 generateWire。
func generateWireAllService(ctx context.Context, projectRootPath string) error {
	appDir := filepath.Join(projectRootPath, "app")
	var entries []os.DirEntry
	var err error
	entries, err = os.ReadDir(appDir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to read app directory: %s\033[m\n", err.Error())
		return err
	}

	var lastErr error
	var processed int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		serviceName := entry.Name()
		servicePath := filepath.Join(appDir, serviceName, "service", "cmd", "server")
		if _, statErr := os.Stat(servicePath); statErr != nil {
			// 没有 service 子目录则跳过
			continue
		}
		processed++
		if genErr := generateWire(ctx, servicePath); genErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: generate for service %s failed: %v\033[m\n", serviceName, genErr)
			lastErr = genErr
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "\033[32mSUCCESS: wire generated for service '%s'\033[m\n", serviceName)
		}
	}
	if processed == 0 {
		err = fmt.Errorf("no services found under %s", appDir)
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return err
	}

	return lastErr
}

// generateWire 在指定目录执行 wire 命令，优先使用全局安装的 wire 可执行文件，如果未找到则回退到使用 go run 运行 wire 命令。
func generateWire(ctx context.Context, serviceRootPath string) error {
	// 尝试查找全局 wire 可执行文件
	wirePath, lookErr := exec.LookPath("wire")

	var cmd *exec.Cmd
	if lookErr == nil {
		// 使用全局 wire
		cmd = exec.CommandContext(ctx, wirePath)
	} else {
		// 回退到使用 go run 运行 wire 命令
		cmd = exec.CommandContext(ctx, "go", "run", "-mod=mod", "github.com/google/wire/cmd/wire")
	}

	cmd.Dir = serviceRootPath
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if runErr := cmd.Run(); runErr != nil {
		if lookErr == nil {
			return fmt.Errorf("failed to run global wire (%s): %w", wirePath, runErr)
		}
		return fmt.Errorf("failed to run fallback 'go run github.com/google/wire/cmd/wire': %w", runErr)
	}

	return nil
}
