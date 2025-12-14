package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	_ = Generate(context.Background(), GeneratorOptions{
		GenerateMain:    true,
		GenerateServer:  true,
		GenerateService: true,
		GenerateData:    true,

		ProjectModule: "test.com/gowind",
		ServiceName:   "user",

		Servers: []string{"rest", "grpc"},

		OutputPath: ".",
	})
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
	err := writeConfigs("./")
	assert.Nil(t, err)
}
