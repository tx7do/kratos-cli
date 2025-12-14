package service

import "strings"

type (
	// GeneratorOptions 保存代码生成器的选项。
	GeneratorOptions struct {
		OutputPath string

		ProjectModule string
		ProjectName   string
		ServiceName   string

		Servers   []string
		DbClients []string

		GenerateServer   bool
		GenerateService  bool
		GenerateData     bool
		GenerateMain     bool
		GenerateMakefile bool
		GenerateConfigs  bool
	}

	GeneratorOption func(*GeneratorOptions)
)

// HasBFFService checks whether the BFF (Backend For Frontend) service is included in the generator options.
func (o GeneratorOptions) HasBFFService() bool {
	isBff := false
	for _, server := range o.Servers {
		if strings.ToLower(server) == "rest" {
			isBff = true
			break
		}
	}
	return isBff
}
