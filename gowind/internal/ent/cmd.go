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

func RunGenerate(cmd *cobra.Command, args []string) error {
	service := args[0]

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
	if err := tidyCmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to run `go mod tidy`: %s\033[m\n", err.Error())
		return err
	}

	servicePath := filepath.Join(inspector.Root, "app", service, "service")

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
	service := args[0]
	namesArg := strings.Join(args[1:], "")
	namesArg = strings.ReplaceAll(namesArg, " ", "")
	names := strings.Split(namesArg, ",")

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
