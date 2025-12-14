package render

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/tx7do/kratos-cli/gowind/internal/service/render/templates"
)

// WriteServerPackageCode writes the server package code to the specified output path.
func WriteServerPackageCode(outputPath string, data ServerTemplateData) error {
	data.Service = strings.ToLower(data.Service)

	outputPath = filepath.Join(outputPath, "/server/")
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	data.Service = snakeToPascal(data.Service)

	switch data.Type {
	case "grpc":
		outputPath = filepath.Join(outputPath, "/", "grpc"+GoFilePostfix)
		outputPath = filepath.Clean(outputPath)
		return renderTemplate[ServerTemplateData](outputPath, data, "grpc_server_"+data.Project, string(templates.GrpcTemplateServerData))

	case "rest":
		outputPath = filepath.Join(outputPath, "/", "rest"+GoFilePostfix)
		outputPath = filepath.Clean(outputPath)
		return renderTemplate[ServerTemplateData](outputPath, data, "rest_server_"+data.Project, string(templates.RestTemplateServerData))

	default:
		return errors.New("unsupported server type: " + data.Type)
	}
}

// WriteInitWireCode writes the initialization wire code to the specified output path.
func WriteInitWireCode(outputPath string, data InitWireTemplateData) error {
	outputPath = filepath.Join(outputPath, "/", data.Package, "/")
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = filepath.Join(outputPath, "/", "init"+GoFilePostfix)
	outputPath = filepath.Clean(outputPath)

	for i, name := range data.ServiceNames {
		data.ServiceNames[i] = snakeToPascal(name)
	}

	return renderTemplate[InitWireTemplateData](outputPath, data, "init_"+data.Package, string(templates.InitTemplateData))
}

// WriteWireCode writes the wire code to the specified output path.
func WriteWireCode(outputPath string, data WireTemplateData) error {
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = filepath.Join(outputPath, "/", "wire"+GoFilePostfix)
	outputPath = filepath.Clean(outputPath)

	data.Service = strings.ToLower(data.Service)

	return renderTemplate[WireTemplateData](outputPath, data, "wire_"+data.Project, string(templates.WireTemplateData))
}

// WriteMainCode writes the main code to the specified output path.
func WriteMainCode(outputPath string, data MainTemplateData) error {
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = filepath.Join(outputPath, "/", "main"+GoFilePostfix)
	outputPath = filepath.Clean(outputPath)

	data.Service = snakeToPascal(data.Service)

	return renderTemplate[MainTemplateData](outputPath, data, "main_"+data.Project, string(templates.MainTemplateData))
}
