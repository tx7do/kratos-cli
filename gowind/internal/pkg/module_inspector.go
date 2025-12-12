package pkg

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ModuleInspector 保存项目根目录和模块名。
type ModuleInspector struct {
	Root    string // 含有 `go.mod` 的目录
	ModPath string // go module path
}

// NewModuleInspectorFromGo 使用 `go list -m -json` 在 startDir（或当前目录）执行，
// 返回模块路径和项目根目录。若在 startDir 执行失败，会向上遍历父目录重试。
func NewModuleInspectorFromGo(startDir string) (*ModuleInspector, error) {
	if startDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		startDir = wd
	}

	g := NewGoCmd("")
	out, err := g.RunUpwardUntilSucceeds(startDir, "list", "-m", "-json")
	if err != nil {
		return nil, err
	}

	var info struct {
		Path string `json:"Path"`
		Dir  string `json:"Dir"`
	}
	if err = json.Unmarshal(out, &info); err != nil {
		return nil, err
	}

	root := info.Dir
	if root == "" {
		// 兜底：使用 go env GOMOD 推断
		g2 := NewGoCmd("")
		if goModOut, e := g2.Output("env", "GOMOD"); e == nil {
			goMod := strings.TrimSpace(string(goModOut))
			if goMod != "" {
				root = filepath.Dir(goMod)
			}
		}
		if root == "" {
			root = startDir
		}
	}

	return &ModuleInspector{
		Root:    root,
		ModPath: info.Path,
	}, nil
}
