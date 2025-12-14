package service

type (
	GeneratorOptions struct {
		OutputPath string

		ProjectModule string
		ServiceName   string

		Servers []string

		GenerateServer   bool
		GenerateService  bool
		GenerateData     bool
		GenerateMain     bool
		GenerateMakefile bool
		GenerateConfigs  bool
	}

	GeneratorOption func(*GeneratorOptions)
)
