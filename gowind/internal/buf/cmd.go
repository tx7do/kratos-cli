package buf

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tx7do/kratos-cli/gowind/internal/pkg"
)

const (
	defaultBufConfigFile    = "buf.yaml"
	bufLockFile             = "buf.lock"
	defaultBufGenConfigFile = "buf.gen.yaml"
)

func RunGenerate(_ *cobra.Command, args []string) error {
	ctx := context.Background()

	inspector, err := pkg.NewModuleInspectorFromGo("")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return err
	}

	// 先在模块根目录运行 `go mod tidy`
	if err = pkg.GoModTidy(ctx, inspector.Root); err != nil {
		return err
	}

	// 确保 buf 已安装
	ensureBufInstalled(ctx)

	apiPath := filepath.Join(inspector.Root, "api")
	if !isDirExists(apiPath) {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: api directory does not exist: %s\033[m\n", apiPath)
		return fmt.Errorf("api directory does not exist: %s", apiPath)
	}

	return GenerateFromPath(ctx, apiPath)
}

// GenerateFromPath 从指定路径生成 Protobuf 代码。
func GenerateFromPath(ctx context.Context, apiPath string) error {
	// 确保 buf 已安装
	ensureBufInstalled(ctx)

	if !isDirExists(apiPath) {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: api directory does not exist: %s\033[m\n", apiPath)
		return fmt.Errorf("api directory does not exist: %s", apiPath)
	}

	if !isBufConfigExists(apiPath) {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: buf config file (%s) does not exist in api directory: %s\033[m\n", defaultBufConfigFile, apiPath)
		return fmt.Errorf("buf config file (%s) does not exist in api directory: %s", defaultBufConfigFile, apiPath)
	}

	if !isBufLockExists(apiPath) {
		fmt.Printf("buf.lock not found in api directory, running `buf dep update`...\n")
		if err := RunBufDepUpdate(ctx, apiPath); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
			return err
		}
	}

	fmt.Printf("Generating proto code from YAML files in api directory...\n")

	yamlFiles, err := scanYAMLFiles(apiPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: failed to scan YAML files: %s\033[m\n", err.Error())
		return err
	}

	if len(yamlFiles) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "\033[33mWARNING: no YAML files found in api directory: %s\033[m\n", apiPath)
		return nil
	}

	fmt.Printf("Running `buf generate` in api directory...\n")
	for _, yamlFile := range yamlFiles {
		fmt.Printf("Using template file: %s\n", yamlFile)
		if err = RunBufGenerate(ctx, apiPath, yamlFile); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
			return err
		}
	}

	fmt.Printf("Protobuf Code generation completed successfully.\n")
	return nil
}

// RunBufDepUpdate 在指定目录执行 `buf dep update`，并将输出转发到标准输出/错误。
func RunBufDepUpdate(ctx context.Context, apiPath string) error {
	cmd := exec.CommandContext(ctx, "buf", "dep", "update")
	cmd.Dir = apiPath
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run `buf dep update`: %w", err)
	}
	return nil
}

// RunBufGenerate 在指定目录执行 `buf generate --template <template>`。
// 如果 template 为空，使用常量 `defaultBufGenConfigFile`。
// 在执行前会调用 ensureBufInstalled 确保 buf 可用，命令输出会直接写入标准输出/错误。
func RunBufGenerate(ctx context.Context, apiPath, template string) error {
	if template == "" {
		template = defaultBufGenConfigFile
	}

	cmd := exec.CommandContext(ctx, "buf", "generate", "--template", template)
	cmd.Dir = apiPath
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run `buf generate --template %s`: %w", template, err)
	}
	return nil
}

// scanYAMLFiles 仅扫描以 `.gen.yaml` 结尾的文件（忽略大小写），返回匹配的文件路径列表。
func scanYAMLFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// 访问失败时继续遍历其他路径
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".gen.yaml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// isDirExists 检查目录是否存在。
func isDirExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

// checkBufInstalled 执行 `buf --version` 来探测 buf 是否安装。
// 返回版本字符串（例: "buf v1.0.0"）或执行失败的错误。
func checkBufInstalled(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "buf", "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run `buf --version`: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return strings.TrimSpace(string(out)), nil
}

// ensureBufInstalled 检查 buf 是否已安装，若未安装则尝试通过 go install 安装。
func ensureBufInstalled(ctx context.Context) {
	if ver, err := checkBufInstalled(ctx); err == nil {
		fmt.Printf("buf is already installed: %s\n", ver)
		return
	}

	fmt.Println("Installing buf...")

	_ = pkg.GoInstall("github.com/bufbuild/buf/cmd/buf@latest")
}

// isBufLockExists 检查 apiPath 下是否存在 `buf.lock` 文件。
func isBufLockExists(apiPath string) bool {
	path := filepath.Join(apiPath, bufLockFile)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// isBufConfigExists 检查 apiPath 下是否存在 `buf.yaml` 文件。
func isBufConfigExists(apiPath string) bool {
	path := filepath.Join(apiPath, defaultBufConfigFile)
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
