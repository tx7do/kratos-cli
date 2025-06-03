package sqlkratos

import (
	"errors"
	"strings"

	"github.com/tx7do/kratos-cli/sql-kratos/internal/render"
)

type DataField render.DataField
type DataFieldArray []render.DataField

func WriteDataPackageCode(
	outputPath string,
	orm string,
	projectName string,
	serviceName string,
	name string,
	moduleName, moduleVersion string,
	protoFields []DataField,
) error {
	var copyDataFields render.DataFieldArray
	for _, field := range protoFields {
		if field.Type == "" {
			continue
		}

		copyDataField := render.DataField{
			Name:    field.Name,
			Type:    field.Type,
			Comment: field.Comment,
		}
		copyDataFields = append(copyDataFields, copyDataField)
	}

	data := render.DataTemplateData{
		Project: projectName,

		Service: serviceName,
		Name:    name,

		Module:  moduleName,
		Version: moduleVersion,

		Fields: copyDataFields,
	}
	switch strings.TrimSpace(strings.ToLower(orm)) {
	case "ent":
		return render.WriteEntDataPackageCode(outputPath, data)

	case "gorm":
		return render.WriteGormDataPackageCode(outputPath, data)

	default:
		return errors.New("sqlproto: unsupported orm: " + orm)
	}
}

func WriteServicePackageCode(
	outputPath string,
	projectName string,
	serviceName string,
	name string,
	targetModuleName, sourceModuleName, moduleVersion string,
	userRepo, isGrpcService bool,
) error {
	data := render.ServiceTemplateData{
		Project: projectName,

		Service: serviceName,
		Name:    name,

		SourceModuleName: sourceModuleName,
		TargetModuleName: targetModuleName,
		Version:          moduleVersion,

		UseRepo: userRepo,
		IsGrpc:  isGrpcService,
	}
	return render.WriteServicePackageCode(outputPath, data)
}

func WriteServerPackageCode(
	outputPath string,
	projectName string,
	serviceType string,
	serviceName string,
	services map[string]string,
) error {
	data := render.ServerTemplateData{
		Project:  projectName,
		Type:     serviceType,
		Service:  serviceName,
		Services: services,
	}
	return render.WriteServerPackageCode(outputPath, data)
}

func WriteInitWireCode(
	outputPath string,

	projectName string,
	postfix string,
	services []string,
) error {
	data := render.InitWireTemplateData{
		Package:      projectName,
		Postfix:      postfix,
		ServiceNames: services,
	}
	return render.WriteInitWireCode(outputPath, data)
}

func WriteWireCode(
	outputPath string,

	projectName string,
	serviceName string,
) error {
	data := render.WireTemplateData{
		Project: projectName,
		Service: serviceName,
	}
	return render.WriteWireCode(outputPath, data)
}

func WriteMainCode(
	outputPath string,

	projectName string,
	serviceName string,

	servers []string,
) error {
	data := render.MainTemplateData{
		Project: projectName,
		Service: serviceName,
		Servers: servers,
	}
	return render.WriteMainCode(outputPath, data)
}
