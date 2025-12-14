package service

import (
	"context"
	"fmt"
	"os"
	"path"

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

var serviceName string

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
		GenerateMain:   true,
		GenerateServer: true,

		ProjectModule: inspector.ModPath,
		ServiceName:   serviceName,

		OutputPath: inspector.Root,
	})
}
