package render

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/tx7do/kratos-cli/sql-kratos/internal/render/templates"
)

// WriteEntDataPackageCode writes the Ent data package code to the specified output path.
func WriteEntDataPackageCode(outputPath string, data DataTemplateData) error {
	outputPath = outputPath + "/data/"
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	return renderTemplate[DataTemplateData](outputPath, data, "ent_data", string(templates.EntRepoTemplate))
}

// WriteGormDataPackageCode writes the GORM data package code to the specified output path.
func WriteGormDataPackageCode(outputPath string, data DataTemplateData) error {
	outputPath = outputPath + "/data/"
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	return renderTemplate[DataTemplateData](outputPath, data, "gorm_data", string(templates.GormRepoTemplate))
}

// WriteServicePackageCode writes the GRPC service package code to the specified output path.
func WriteServicePackageCode(outputPath string, data ServiceTemplateData) error {
	outputPath = outputPath + "/service/"
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	return renderTemplate[ServiceTemplateData](outputPath, data, "service", string(templates.ServiceTemplate))
}

// WriteServerPackageCode writes the server package code to the specified output path.
func WriteServerPackageCode(outputPath string, data ServerTemplateData) error {
	data.Service = strings.ToLower(data.Service)

	outputPath = outputPath + "/server/"
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	data.Service = snakeToPascal(data.Service)

	switch data.Type {
	case "grpc":
		outputPath = outputPath + "/" + "grpc" + GoFilePostfix
		outputPath = filepath.Clean(outputPath)
		return renderTemplate[ServerTemplateData](outputPath, data, "grpc_server_"+data.Project, string(templates.GrpcServerTemplate))

	case "rest":
		outputPath = outputPath + "/" + "rest" + GoFilePostfix
		outputPath = filepath.Clean(outputPath)
		return renderTemplate[ServerTemplateData](outputPath, data, "rest_server_"+data.Project, string(templates.RestServerTemplate))

	default:
		return errors.New("unsupported server type: " + data.Type)
	}
}

// WriteInitWireCode writes the initialization wire code to the specified output path.
func WriteInitWireCode(outputPath string, data InitWireTemplateData) error {
	outputPath = outputPath + "/" + data.Package + "/"
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = outputPath + "/" + "init" + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	for i, name := range data.ServiceNames {
		data.ServiceNames[i] = snakeToPascal(name)
	}

	return renderTemplate[InitWireTemplateData](outputPath, data, "init_"+data.Package, string(templates.InitTemplate))
}

// WriteWireCode writes the wire code to the specified output path.
func WriteWireCode(outputPath string, data WireTemplateData) error {
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = outputPath + "/" + "wire" + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	data.Service = strings.ToLower(data.Service)

	return renderTemplate[WireTemplateData](outputPath, data, "wire_"+data.Project, string(templates.WireTemplate))
}

// WriteMainCode writes the main code to the specified output path.
func WriteMainCode(outputPath string, data MainTemplateData) error {
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = outputPath + "/" + "main" + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	data.Service = snakeToPascal(data.Service)

	return renderTemplate[MainTemplateData](outputPath, data, "main_"+data.Project, string(templates.MainTemplate))
}
