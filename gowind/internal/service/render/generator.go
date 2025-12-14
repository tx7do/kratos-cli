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

	switch strings.ToLower(data.Type) {
	case "grpc":
		outputPath = filepath.Join(outputPath, "/", "grpc"+GoFilePostfix)
		outputPath = filepath.Clean(outputPath)
		return renderTemplate[ServerTemplateData](outputPath, data, "grpc_server_"+data.Project, string(templates.GrpcServerTemplate))

	case "rest":
		outputPath = filepath.Join(outputPath, "/", "rest"+GoFilePostfix)
		outputPath = filepath.Clean(outputPath)
		return renderTemplate[ServerTemplateData](outputPath, data, "rest_server_"+data.Project, string(templates.RestServerTemplate))

	default:
		return errors.New("unsupported server type: " + data.Type)
	}
}

// WriteDataPackageCode writes the data package code to the specified output path.
func WriteDataPackageCode(outputPath string, data DataTemplateData) error {
	data.Service = strings.ToLower(data.Service)

	outputPath = filepath.Join(outputPath, "/data/")
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPathData := filepath.Join(outputPath, "/", "data"+GoFilePostfix)
	outputPathData = filepath.Clean(outputPathData)
	if err := renderTemplate[DataTemplateData](outputPathData, data, "data"+data.Project, string(templates.DataTemplate)); err != nil {
		return err
	}

	for _, dbc := range data.DBClients {
		switch strings.ToLower(dbc) {
		case "gorm":
			outputPathGorm := filepath.Join(outputPath, "/", "gorm_client"+GoFilePostfix)
			outputPathGorm = filepath.Clean(outputPathGorm)
			if err := renderTemplate[DataTemplateData](outputPathGorm, data, "gorm_client"+data.Project, string(templates.GormClientTemplate)); err != nil {
				return err
			}

			gormDir := filepath.Join(outputPath, "/gorm/")
			if err := os.MkdirAll(gormDir, os.ModePerm); err != nil {
				return err
			}

			gormModelsDir := filepath.Join(gormDir, "/models/")
			if err := os.MkdirAll(gormModelsDir, os.ModePerm); err != nil {
				return err
			}

			outputPathGormInit := filepath.Join(gormDir, "init"+GoFilePostfix)
			outputPathGormInit = filepath.Clean(outputPathGormInit)
			if err := renderTemplate[DataTemplateData](outputPathGormInit, data, "gorm_init"+data.Project, string(templates.GormInitTemplate)); err != nil {
				return err
			}

		case "ent", "entgo":
			outputPathEnt := filepath.Join(outputPath, "/", "ent_client"+GoFilePostfix)
			outputPathEnt = filepath.Clean(outputPathEnt)
			if err := renderTemplate[DataTemplateData](outputPathEnt, data, "ent_client"+data.Project, string(templates.EntClientTemplate)); err != nil {
				return err
			}

			entDir := filepath.Join(outputPath, "/ent/")
			if err := os.MkdirAll(entDir, os.ModePerm); err != nil {
				return err
			}

			entSchemaDir := filepath.Join(entDir, "/schema/")
			if err := os.MkdirAll(entSchemaDir, os.ModePerm); err != nil {
				return err
			}

		case "redis":
			outputPathRedis := filepath.Join(outputPath, "/", "redis_client"+GoFilePostfix)
			outputPathRedis = filepath.Clean(outputPathRedis)
			if err := renderTemplate[DataTemplateData](outputPathRedis, data, "redis_client"+data.Project, string(templates.RedisClientTemplate)); err != nil {
				return err
			}

		default:
			return errors.New("unsupported data client: " + dbc)
		}
	}

	return nil
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

	var functionData InitWireFunctionTemplateData
	functionData.Package = data.Package
	for _, name := range data.ServiceNames {
		functionData.Functions = append(functionData.Functions, "New"+snakeToPascal(name)+data.Postfix)
	}

	return renderTemplate[InitWireFunctionTemplateData](outputPath, functionData, "init_"+data.Package, string(templates.InitTemplate))
}

// WriteInitWireFunctionCode writes the initialization wire code to the specified output path.
func WriteInitWireFunctionCode(outputPath string, data InitWireFunctionTemplateData) error {
	outputPath = filepath.Join(outputPath, "/", data.Package, "/")
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = filepath.Join(outputPath, "/", "init"+GoFilePostfix)
	outputPath = filepath.Clean(outputPath)

	return renderTemplate[InitWireFunctionTemplateData](outputPath, data, "init_"+data.Package, string(templates.InitTemplate))
}

// WriteWireCode writes the wire code to the specified output path.
func WriteWireCode(outputPath string, data WireTemplateData) error {
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = filepath.Join(outputPath, "/", "wire"+GoFilePostfix)
	outputPath = filepath.Clean(outputPath)

	data.Service = strings.ToLower(data.Service)

	return renderTemplate[WireTemplateData](outputPath, data, "wire_"+data.Project, string(templates.WireTemplate))
}

// WriteMainCode writes the main code to the specified output path.
func WriteMainCode(outputPath string, data MainTemplateData) error {
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = filepath.Join(outputPath, "/", "main"+GoFilePostfix)
	outputPath = filepath.Clean(outputPath)

	data.Service = snakeToPascal(data.Service)

	return renderTemplate[MainTemplateData](outputPath, data, "main_"+data.Project, string(templates.MainTemplate))
}
