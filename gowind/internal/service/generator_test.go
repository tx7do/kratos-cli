package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	err := Generate(context.Background(), GeneratorOptions{
		GenerateMain:     true,
		GenerateServer:   true,
		GenerateService:  true,
		GenerateData:     true,
		GenerateMakefile: true,
		GenerateConfigs:  true,

		ProjectModule: "github.com/gowind-example",
		ProjectName:   "gowind-example",
		ServiceName:   "user",

		Servers:   []string{"rest", "grpc"},
		DbClients: []string{"ent", "redis"},

		OutputPath: "./test",
	})
	assert.Nil(t, err)
}

func TestAppendServiceName(t *testing.T) {
	//names := []string{"user", "order"}
	err := appendServiceName("./", "test", "user", false)
	assert.Nil(t, err)

	err = appendServiceName("./", "test", "order", false)
	assert.Nil(t, err)

	err = appendServiceName("./", "test", "admin", true)
	assert.Nil(t, err)

	err = appendServiceName("./", "test", "front", true)
	assert.Nil(t, err)
}

func TestWriteMakefile(t *testing.T) {
	err := writeMakefile("./")
	assert.Nil(t, err)
}

func TestWriteConfigs(t *testing.T) {
	err := writeConfigs("./configs")
	assert.Nil(t, err)
}

func TestExtractProjectName(t *testing.T) {
	projectModule := "github.com/gowind-example"
	projectName := extractProjectName(projectModule)
	assert.Equal(t, "gowind-example", projectName)

	projectModule = "gowind-example"
	projectName = extractProjectName(projectModule)
	assert.Equal(t, "gowind-example", projectName)
}
