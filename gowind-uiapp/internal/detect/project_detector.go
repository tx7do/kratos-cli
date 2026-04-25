package detect

import (
	"os"
	"path/filepath"
)

type ProjectInfo struct {
	Root         string   `json:"Root"`
	GoVersion    string   `json:"GoVersion"`
	ModPath      string   `json:"ModPath"`
	Main         bool     `json:"Main"`
	Version      string   `json:"Version"`
	Replace      *Module  `json:"Replace,omitempty"`
	Dependencies []Module `json:"Dependencies,omitempty"`

	Services []string `json:"Services,omitempty"`
	HasApi   bool     `json:"HasApi,omitempty"`
}

type ProjectDetector struct {
}

func NewProjectDetector() *ProjectDetector {
	return &ProjectDetector{}
}

// Detect 检测指定路径下的 Go 项目，返回项目的基本信息。
func (pd *ProjectDetector) Detect(projectPath string) (*ProjectInfo, error) {
	inspector, err := NewModuleInspectorFromGo(projectPath)
	if err != nil {
		return nil, err
	}

	var pi ProjectInfo
	pi.Root = projectPath
	pi.ModPath = inspector.ModPath
	pi.GoVersion = inspector.GoVersion
	pi.Main = inspector.Main
	pi.Version = inspector.Version
	pi.Replace = inspector.Replace
	pi.Dependencies = inspector.Dependencies

	// 收集服务列表
	services, err := pd.collectServices(projectPath)
	if err == nil {
		pi.Services = services
	}

	// 检测是否包含 API 定义
	pi.HasApi = pd.collectApi(projectPath)

	return &pi, nil
}

// collectServices 收集项目中的服务列表，假设服务目录位于 projectPath/app 下。
func (pd *ProjectDetector) collectServices(projectPath string) ([]string, error) {
	var services []string

	servicesPath := filepath.Join(projectPath, "app")
	entries, err := os.ReadDir(servicesPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			services = append(services, entry.Name())
		}
	}

	return services, nil
}

// collectApi 检测项目中是否包含 API 定义，假设 API 定义位于 projectPath/api 下。
func (pd *ProjectDetector) collectApi(projectPath string) bool {
	apiPath := filepath.Join(projectPath, "api")
	_, err := os.Stat(apiPath)
	if err != nil {
		return false
	}

	bufYamlPath := filepath.Join(apiPath, "buf.yaml")
	_, err = os.Stat(bufYamlPath)
	if err != nil {
		return false
	}

	protosPath := filepath.Join(apiPath, "protos")
	_, err = os.Stat(protosPath)
	if err != nil {
		return false
	}

	return true
}
