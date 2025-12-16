package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// GoCmd 封装 go 命令执行。
type GoCmd struct {
	Dir     string        // 工作目录
	Env     []string      // 额外环境变量（追加到 os.Environ()）
	Timeout time.Duration // 若 >0，为每次执行设置超时

	Stdin  io.Reader
	Stdout io.Writer // 若为 nil，Run 会使用 os.Stdout
	Stderr io.Writer // 若为 nil，Run 会使用 os.Stderr
}

// NewGoCmd 返回一个默认的 GoCmd（无超时）。
func NewGoCmd(dir string) *GoCmd {
	return &GoCmd{Dir: dir}
}

// NewGoCmdWithTimeout 返回一个带超时的 GoCmd。
func NewGoCmdWithTimeout(dir string, timeout time.Duration) *GoCmd {
	return &GoCmd{Dir: dir, Timeout: timeout}
}

// prepareEnv 返回完整环境变量切片。
func (g *GoCmd) prepareEnv() []string {
	env := os.Environ()
	if len(g.Env) > 0 {
		env = append(env, g.Env...)
	}
	return env
}

// makeContext 返回基于 g.Timeout 的 context。
func (g *GoCmd) makeContext() (context.Context, context.CancelFunc) {
	if g.Timeout > 0 {
		return context.WithTimeout(context.Background(), g.Timeout)
	}
	return context.WithCancel(context.Background())
}

// Run 直接执行 go 命令，输出到 GoCmd 指定的 Stdout/Stderr（或终端）。
func (g *GoCmd) Run(args ...string) error {
	fmt.Printf("go %s\n", joinArgs(args))

	ctx, cancel := g.makeContext()
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", args...)
	if g.Dir != "" {
		cmd.Dir = g.Dir
	}
	cmd.Env = g.prepareEnv()
	if g.Stdin != nil {
		cmd.Stdin = g.Stdin
	}
	if g.Stdout != nil {
		cmd.Stdout = g.Stdout
	} else {
		cmd.Stdout = os.Stdout
	}
	if g.Stderr != nil {
		cmd.Stderr = g.Stderr
	} else {
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

// Output 执行并返回 stdout（不包含 stderr）。
func (g *GoCmd) Output(args ...string) ([]byte, error) {
	ctx, cancel := g.makeContext()
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", args...)
	if g.Dir != "" {
		cmd.Dir = g.Dir
	}
	cmd.Env = g.prepareEnv()
	if g.Stdin != nil {
		cmd.Stdin = g.Stdin
	}
	return cmd.Output()
}

// CombinedOutput 执行并返回 stdout+stderr。
func (g *GoCmd) CombinedOutput(args ...string) ([]byte, error) {
	ctx, cancel := g.makeContext()
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", args...)
	if g.Dir != "" {
		cmd.Dir = g.Dir
	}
	cmd.Env = g.prepareEnv()
	if g.Stdin != nil {
		cmd.Stdin = g.Stdin
	}
	return cmd.CombinedOutput()
}

// RunUpwardUntilSucceeds 从 startDir 开始向上遍历父目录，尝试执行 go 命令，
// 直到某一目录执行成功或到达根目录。返回最后一次成功的输出（combined）。
func (g *GoCmd) RunUpwardUntilSucceeds(startDir string, args ...string) ([]byte, error) {
	dir := startDir
	if dir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = wd
	}
	for {
		g.Dir = dir
		out, err := g.CombinedOutput(args...)
		if err == nil {
			return out, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// 到根仍失败，返回最后一个错误和输出
			return out, err
		}
		dir = parent
	}
}

// Helper: join args for logging
func joinArgs(args []string) string {
	var buf bytes.Buffer
	for i, a := range args {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(a)
	}
	return buf.String()
}

// GoInstall 使用 `go install` 安装指定路径的包。
// 若路径中不包含版本号，则默认使用 @latest。
func GoInstall(paths ...string) error {
	for _, p := range paths {
		if !containsAt(p) {
			p += "@latest"
		}
		g := NewGoCmd("")
		if err := g.Run("install", p); err != nil {
			return err
		}
	}
	return nil
}

func containsAt(s string) bool {
	for _, c := range s {
		if c == '@' {
			return true
		}
	}
	return false
}
