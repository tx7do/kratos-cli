package service

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/tx7do/kratos-cli/gowind/internal/pkg"
)

// CmdService represents the service command
var CmdService = &cobra.Command{
	Use:     "service [name]",
	Aliases: []string{"svc"},
	Short:   "create a new service scaffold",
	Long:    "Create a new microservice inside the current workspace. Example: gow new service user",
	Args:    cobra.ExactArgs(1),
	Run:     run,
}

var (
	serviceName string
	Servers     []string
	DbClients   []string
)

func init() {
	Servers = []string{"grpc"}
	DbClients = []string{"ent"}

	CmdService.Flags().StringArrayVarP(&Servers, "servers", "s", []string{"grpc"}, "Specify which server types to generate (grpc, rest, asynq, sse...)")
	CmdService.Flags().StringArrayVarP(&DbClients, "db-clients", "d", []string{"ent"}, "Specify which database clients to generate (gorm, ent, redis, clickhouse...)")
}

func extractProjectName(module string) string {
	module = strings.TrimSpace(module)
	if module == "" {
		return ""
	}

	if strings.Contains(module, "/") {
		parts := strings.Split(module, "/")
		for i := len(parts) - 1; i >= 0; i-- {
			seg := strings.TrimSpace(parts[i])
			if seg != "" {
				return seg
			}
		}
	}

	return module
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is service name?",
			Help:    "Created service name.",
		}
		err := survey.AskOne(prompt, &serviceName)
		if err != nil || serviceName == "" {
			return
		}
	} else {
		serviceName = args[0]
	}

	inspector, err := pkg.NewModuleInspectorFromGo("")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err.Error())
		return
	}

	servicePath := path.Join(inspector.Root, "/app/", serviceName, "/service")

	if pkg.IsDirExists(servicePath) {
		_, _ = fmt.Fprintf(os.Stderr, "\033[31mERROR: Service directory %s already exists\033[m\n", servicePath)
		return
	}

	_ = Generate(context.Background(), GeneratorOptions{
		GenerateMain:     true,
		GenerateServer:   true,
		GenerateService:  true,
		GenerateData:     true,
		GenerateMakefile: true,
		GenerateConfigs:  true,

		ProjectName:   extractProjectName(inspector.ModPath),
		ProjectModule: inspector.ModPath,
		ServiceName:   serviceName,

		Servers:   Servers,
		DbClients: DbClients,

		OutputPath: inspector.Root,
	})
}
