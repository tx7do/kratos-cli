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
	inspector, err := pkg.NewModuleInspectorFromGo("")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return err
	}

	// 先在模块根目录运行 `go mod tidy`
	tidyCmd := exec.CommandContext(context.Background(), "go", "mod", "tidy")
	tidyCmd.Dir = inspector.Root
	tidyCmd.Env = os.Environ()
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	if err = tidyCmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to run `go mod tidy`: %s\033[m\n", err.Error())
		return err
	}

	// 如果没有传入 service 或者为空，则遍历所有 service
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		appDir := filepath.Join(inspector.Root, "app")
		entries, err := os.ReadDir(appDir)
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
			svc := entry.Name()
			servicePath := filepath.Join(appDir, svc, "service", "cmd", "server")
			if _, statErr := os.Stat(servicePath); statErr != nil {
				// 没有 service 子目录则跳过
				continue
			}
			processed++
			if genErr := generateWire(servicePath); genErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: generate for service %s failed: %v\033[m\n", svc, genErr)
				lastErr = genErr
			} else {
				fmt.Printf("Generated wire for service %s successfully.\n", svc)
			}
		}
		if processed == 0 {
			err = fmt.Errorf("no services found under %s", appDir)
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
			return err
		}
		return lastErr
	}

	service := args[0]
	servicePath := filepath.Join(inspector.Root, "app", service, "service", "cmd", "server")
	return generateWire(servicePath)
}

func generateWire(dir string) error {
	ctx := context.Background()

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

	cmd.Dir = dir
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
