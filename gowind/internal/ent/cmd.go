package ent

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

func RunGenerate(_ *cobra.Command, args []string) error {
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
			servicePath := filepath.Join(appDir, svc, "service")
			if _, statErr := os.Stat(servicePath); statErr != nil {
				// 没有 service 子目录则跳过
				continue
			}
			processed++
			if genErr := generateEnt(servicePath); genErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: generate for service %s failed: %v\033[m\n", svc, genErr)
				lastErr = genErr
			} else {
				fmt.Printf("Generated ent for service %s successfully.\n", svc)
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
	servicePath := filepath.Join(inspector.Root, "app", service, "service")
	return generateEnt(servicePath)
}

func generateEnt(servicePath string) error {
	target := filepath.Join(servicePath, "internal", "data", "ent", "schema")
	e := NewEntCmd(target)
	return e.RunGenerate(
		"--feature", "privacy",
		"--feature", "entql",
		"--feature", "sql/modifier",
		"--feature", "sql/upsert",
		"--feature", "sql/lock",
	)
}

func RunAdd(cmd *cobra.Command, args []string) error {
	// 最少需要 service 和 schemas
	if len(args) < 2 {
		_ = cmd.Help()
		return fmt.Errorf("usage: ent add <service> <schemas>")
	}

	service := strings.TrimSpace(args[0])
	if service == "" {
		_ = cmd.Help()
		return fmt.Errorf("service name is empty")
	}

	// 支持多个参数或一个逗号分隔字符串，去掉空格并拆分
	namesArg := strings.Join(args[1:], "")
	namesArg = strings.ReplaceAll(namesArg, " ", "")
	if namesArg == "" {
		_ = cmd.Help()
		return fmt.Errorf("no schema names provided")
	}

	names := strings.Split(namesArg, ",")
	// 过滤空名称
	var filtered []string
	for _, n := range names {
		n = strings.TrimSpace(n)
		if n != "" {
			filtered = append(filtered, n)
		}
	}
	if len(filtered) == 0 {
		_ = cmd.Help()
		return fmt.Errorf("no valid schema names after parsing")
	}

	inspector, err := pkg.NewModuleInspectorFromGo("")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return err
	}

	servicePath := filepath.Join(inspector.Root, "app", service, "service")
	target := filepath.Join(servicePath, "internal", "data")

	e := NewEntCmd(target)
	return e.RunNew(names)
}
