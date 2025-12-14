package service

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tx7do/go-utils/stringcase"
	"github.com/tx7do/kratos-cli/gowind/internal/pkg"
)

func Generate(ctx context.Context, opts GeneratorOptions) error {
	var err error

	// 生成server层代码
	if opts.GenerateServer {
		serverPackagePath := filepath.Join(opts.OutputPath, "/app/", opts.ServiceName, "/service/internal/")
		if err = generateServerPackageCode(
			serverPackagePath,
			opts.ProjectModule,
			opts.ServiceName,
			opts.Servers,
		); err != nil {
			return err
		}
	}

	// 生成service层代码
	if opts.GenerateService {
		servicePackagePath := filepath.Join(opts.OutputPath, "/app/", opts.ServiceName, "/service/internal/")
		if err = generateServicePackageCode(
			servicePackagePath,
			opts.ProjectModule,
			opts.ServiceName,
			[]string{},
		); err != nil {
			return err
		}
	}

	// 生成data层代码
	if opts.GenerateData {
		dataPackagePath := filepath.Join(opts.OutputPath, "/app/", opts.ServiceName, "/service/internal/")
		if err = generateDataPackageCode(
			dataPackagePath,
			opts.ProjectModule,
			opts.ServiceName,
			[]string{},
		); err != nil {
			return err
		}
	}

	// 生成main包代码
	if opts.GenerateMain {
		mainPackagePath := filepath.Join(opts.OutputPath, "/app/", opts.ServiceName, "/service/cmd/server")
		if err = generateMainPackageCode(
			mainPackagePath,
			opts.ProjectModule,
			opts.ServiceName,
			opts.Servers,
		); err != nil {
			return err
		}
	}

	// 生成Makefile
	if opts.GenerateMakefile {
		makefilePath := filepath.Join(opts.OutputPath, "/app/", opts.ServiceName, "/service")
		if err = writeMakefile(makefilePath); err != nil {
			return err
		}
	}

	// 生成configs
	if opts.GenerateConfigs {
		configsPath := filepath.Join(opts.OutputPath, "/app/", opts.ServiceName, "/service/configs")
		if err = writeConfigs(configsPath); err != nil {
			return err
		}
	}

	return nil
}

func generateServerPackageCode(
	outputPath string,
	projectName string,
	serviceName string,
	servers []string,
) error {

	for _, server := range servers {
		if err := WriteServerPackageCode(
			outputPath,
			projectName, server, serviceName,
		); err != nil {
			return err
		}
	}

	return WriteInitWireCode(outputPath, "server", "Server", servers)
}

func generateServicePackageCode(
	outputPath string,
	projectName string,
	serviceName string,
	services []string,
) error {
	return WriteInitWireCode(outputPath, "service", "Service", services)
}

func generateDataPackageCode(
	outputPath string,
	projectName string,
	serviceName string,
	repos []string,
) error {
	return WriteInitWireCode(outputPath, "data", "Repo", repos)
}

func generateMainPackageCode(
	outputPath string,
	projectName string, serviceName string,
	servers []string,
) error {
	if err := WriteMainCode(
		outputPath,
		projectName, serviceName,
		servers,
	); err != nil {
		return err
	}

	return WriteWireCode(
		outputPath,
		projectName, serviceName,
	)
}

// writeMakefile 生成默认的 Makefile 到指定目录。
func writeMakefile(outputPath string) error {
	outputPath = outputPath + "/"
	outputPath = filepath.Clean(outputPath)
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	makefilePath := path.Join(outputPath, "Makefile")
	makefilePath = filepath.Clean(makefilePath)

	const makefileContent string = `
include ../../../app.mk
`
	if !pkg.IsFileExists(makefilePath) {
		if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
			return fmt.Errorf("write Makefile: %w", err)
		}
	}

	return nil
}

// writeConfigs 生成默认的配置文件到指定目录。
func writeConfigs(outputPath string) error {
	const dataYaml string = `
data:
  database:
    driver: "postgres"
    source: "host=postgres port=5432 user=postgres password=<you_password> dbname=gwa sslmode=disable"
#    driver: "mysql"
#    source: "root:<you_password>@tcp(localhost:3306)/gwa?parseTime=true&charset=utf8mb4&loc=Asia%2FShanghai"
    migrate: true
    debug: false
    enable_trace: false
    enable_metrics: false
    max_idle_connections: 25
    max_open_connections: 25
    connection_max_lifetime: 300s

  redis:
    addr: "redis:6379"
    password: "<you_password>"
    dial_timeout: 10s
    read_timeout: 0.4s
    write_timeout: 0.6s
`

	const serverYaml string = `
server:
  grpc:
    addr: ":8899"
    timeout: 10s
    middleware:
`

	const loggerYaml string = `
logger:
  type: std # Options: std, file, fluent, zap, logrus, aliyun, tencent

  fluent:
    endpoint: "tcp://localhost:24224"

  zap:
    level: "debug"
    filename: "./logs/info.log"
    max_size: 1
    max_age: 30
    max_backups: 5

  logrus:
    level: "debug"
    formatter: "text"
    timestamp_format: "2006-01-02 15:04:05"
    disable_colors: false
    disable_timestamp: false

  aliyun:
    endpoint: ""
    project: ""
    access_key: "<access_key>"
    access_secret: "<access_secret>"

  tencent:
    endpoint: ""
    topic_id:
    access_key: "<access_key>"
    access_secret: "<access_secret>"
`

	const clientYaml string = `
client:
  grpc:
    timeout: 10s
    middleware:
      enable_logging: true
      enable_recovery: true
      enable_tracing: true
      enable_validate: true
      enable_circuit_breaker: true
      enable_metadata: true
      auth:
        method: "HS256"
        key: "some_api_key"
`

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return fmt.Errorf("create configs dir: %w", err)
	}

	serverYamlFile := filepath.Join(outputPath, "server.yaml")
	if !pkg.IsFileExists(serverYamlFile) {
		if err := os.WriteFile(serverYamlFile, []byte(serverYaml), 0644); err != nil {
			return fmt.Errorf("write server.yaml file: %w", err)
		}
	}

	dataYamlFile := filepath.Join(outputPath, "data.yaml")
	if !pkg.IsFileExists(dataYamlFile) {
		if err := os.WriteFile(dataYamlFile, []byte(dataYaml), 0644); err != nil {
			return fmt.Errorf("write data.yaml file: %w", err)
		}
	}

	loggerYamlFile := filepath.Join(outputPath, "logger.yaml")
	if !pkg.IsFileExists(loggerYamlFile) {
		if err := os.WriteFile(loggerYamlFile, []byte(loggerYaml), 0644); err != nil {
			return fmt.Errorf("write logger.yaml file: %w", err)
		}
	}

	clientYamlFile := filepath.Join(outputPath, "client.yaml")
	if !pkg.IsFileExists(clientYamlFile) {
		if err := os.WriteFile(clientYamlFile, []byte(clientYaml), 0644); err != nil {
			return fmt.Errorf("write client.yaml file: %w", err)
		}
	}

	return nil
}

// appendServiceName 向 pkg/service/name.go 文件追加服务名称常量定义。
func appendServiceName(outputPath string, projectName, serviceName string, isBff bool) error {
	pkgDir := filepath.Join(outputPath, "pkg", "service")
	if err := os.MkdirAll(pkgDir, os.ModePerm); err != nil {
		return fmt.Errorf("create pkg/service dir: %w", err)
	}

	servicePostfix := "service"
	if isBff {
		servicePostfix = "bff"
	}

	// 常量名与值
	constName := fmt.Sprintf("%sService", stringcase.UpperCamelCase(serviceName))
	constValue := fmt.Sprintf("%s-%s-%s", stringcase.LowerCamelCase(projectName), strings.ToLower(serviceName), servicePostfix)

	// 行格式，带缩进
	fieldLine := fmt.Sprintf("    %s = %q", constName, constValue)

	serviceNamePath := filepath.Join(pkgDir, "name.go")

	// 文件不存在：创建包含 const 块的初始文件
	if !pkg.IsFileExists(serviceNamePath) {
		content := fmt.Sprintf("package service\n\nconst (\n%s\n)\n", fieldLine)
		if err := os.WriteFile(serviceNamePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("write service name file: %w", err)
		}
		return nil
	}

	// 文件存在：读取并检查是否已包含常量名
	data, err := os.ReadFile(serviceNamePath)
	if err != nil {
		return fmt.Errorf("read service name file: %w", err)
	}
	text := string(data)
	if strings.Contains(text, constName) {
		// 已包含，不需插入
		return nil
	}

	// 找到第一个 const ( ... ) 块并在闭合 ) 之前插入新行
	constIdx := strings.Index(text, "const (")
	if constIdx >= 0 {
		// 在 const ( 后寻找对应的第一个 )（简单实现，适用于代码生成的文件）
		closeIdx := strings.Index(text[constIdx:], ")")
		if closeIdx >= 0 {
			insertPos := constIdx + closeIdx
			newText := text[:insertPos] + "\n" + fieldLine + "\n" + text[insertPos:]
			if err = os.WriteFile(serviceNamePath, []byte(newText), 0644); err != nil {
				return fmt.Errorf("write service name file after insert: %w", err)
			}
			return nil
		}
	}

	// 未找到 const 块，直接在文件末尾追加一个新的 const 块
	appendContent := fmt.Sprintf("\nconst (\n%s\n)\n", fieldLine)
	newText := text + appendContent
	if err := os.WriteFile(serviceNamePath, []byte(newText), 0644); err != nil {
		return fmt.Errorf("append service name file: %w", err)
	}

	return nil
}
