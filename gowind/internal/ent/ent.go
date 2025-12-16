package ent

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/tx7do/kratos-cli/gowind/internal/pkg"
)

// EntCmd 封装 ent 命令调用的数据和行为。
type EntCmd struct {
	Args      []string
	TargetDir string
	Timeout   time.Duration
	goCmd     *pkg.GoCmd
}

func NewEntCmd(targetDir string) *EntCmd {
	if targetDir == "" {
		targetDir = "internal/data/ent/schema"
	}
	timeout := 2 * time.Minute
	e := &EntCmd{
		TargetDir: filepath.Clean(targetDir),
		Timeout:   timeout,
	}
	e.goCmd = pkg.NewGoCmdWithTimeout(targetDir, timeout)
	return e
}

// tryGoRunFirst 尝试在项目中通过 `go run entgo.io/ent/cmd/ent ...` 执行命令（向上查找可执行目录）。
// 成功时打印输出并返回 nil；失败时返回最后一次错误以便调用方决定后续处理。
func (e *EntCmd) tryGoRunFirst(args ...string) error {
	out, err := e.goCmd.RunUpwardUntilSucceeds("", args...)
	if err == nil {
		// 打印 combined output（stdout+stderr）
		if len(out) > 0 {
			fmt.Print(string(out))
		}
		return nil
	}
	return err
}

// runGlobalEntIfAvailable 尝试在 PATH 中查找全局 `ent` 可执行程序并执行（遵循 e.Timeout）。
func (e *EntCmd) runGlobalEntIfAvailable(args ...string) error {
	entPath, err := exec.LookPath("ent")
	if err != nil {
		return fmt.Errorf("no local go-run ent succeeded and no global ent found: %w", err)
	}

	// 准备上下文（带超时）
	var ctx context.Context
	var cancel context.CancelFunc
	if e.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), e.Timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, entPath, args...)

	// 将工作目录设置为 `e.TargetDir` 的绝对路径
	if e.TargetDir != "" {
		if abs, err := filepath.Abs(e.TargetDir); err == nil {
			cmd.Dir = abs
		} else {
			// 无法获取绝对路径时仍可使用原始值
			cmd.Dir = e.TargetDir
		}
	}

	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println("running global ent:", entPath, args)
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("global ent failed: %w", err)
	}
	return nil
}

// RunNew 使用优先使用项目内的 go run，如果失败则尝试全局 ent 可执行程序来执行 `ent new`。
func (e *EntCmd) RunNew(names []string) error {
	if len(names) == 0 {
		return errors.New("at least one schema name is required")
	}
	if err := os.MkdirAll(e.TargetDir, 0o755); err != nil {
		return fmt.Errorf("create schema dir failed: %w", err)
	}

	// 尝试 go run first
	goArgs := []string{"run", "entgo.io/ent/cmd/ent", "new"}
	goArgs = append(goArgs, names...)

	if err := e.tryGoRunFirst(goArgs...); err == nil {
		return nil
	}

	// 如果 go run 不可用，则尝试全局 ent
	entArgs := append([]string{"new"}, names...)
	return e.runGlobalEntIfAvailable(entArgs...)
}

// RunGenerate 使用优先使用项目内的 go run，如果失败则尝试全局 ent 可执行程序来执行 `ent generate`。
// extraArgs 可传入像 `--feature` 这样的额外参数。
func (e *EntCmd) RunGenerate(extraArgs ...string) error {
	targetDir := e.TargetDir
	if targetDir == "" {
		targetDir = "."
	}

	goArgs := []string{"run", "entgo.io/ent/cmd/ent", "generate"}
	if len(extraArgs) > 0 {
		goArgs = append(goArgs, extraArgs...)
	}
	goArgs = append(goArgs, targetDir)

	// 尝试 go run first
	if err := e.tryGoRunFirst(goArgs...); err == nil {
		return nil
	}

	// 如果 go run 不可用，则尝试全局 ent
	entArgs := []string{"generate"}
	if len(extraArgs) > 0 {
		entArgs = append(entArgs, extraArgs...)
	}
	entArgs = append(entArgs, targetDir)
	return e.runGlobalEntIfAvailable(entArgs...)
}
