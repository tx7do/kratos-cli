package sqlproto

import (
	"errors"
	"log"
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/tx7do/kratos-cli/sql-proto/internal/render"
)

type ProtoField render.ProtoField
type ProtoFieldArray []render.ProtoField

func WriteServiceProto(
	outputPath string,
	serviceType string,
	targetModuleName, sourceModuleName, moduleVersion string,
	tableName string,
	tableComment string,
	protoFields ProtoFieldArray,
) error {
	switch strings.TrimSpace(strings.ToLower(serviceType)) {
	case "grpc":
		data := render.GrpcProtoTemplateData{
			Module:  targetModuleName,
			Version: moduleVersion,

			Name:    inflection.Singular(tableName),
			Comment: render.RemoveTableCommentSuffix(tableComment),
			Fields:  render.ProtoFieldArray(protoFields),
		}
		return render.WriteGrpcServiceProto(outputPath, data)

	case "rest":
		data := render.RestProtoTemplateData{
			SourceModule: sourceModuleName,
			TargetModule: targetModuleName,
			Version:      moduleVersion,

			Name:    inflection.Singular(tableName),
			Comment: render.RemoveTableCommentSuffix(tableComment),
		}
		return render.WriteRestServiceProto(outputPath, data)

	default:
		return errors.New("sqlproto: unsupported service type: " + serviceType)
	}
}

func WriteServicesProto(
	outputPath string,
	serviceType string,
	targetModuleName, sourceModuleName, moduleVersion string,
	tables TableDataArray,
) error {
	var protoFields ProtoFieldArray

	for i := 0; i < len(tables); i++ {
		table := tables[i]

		protoFields = make(ProtoFieldArray, 0, len(table.Fields))
		for n := 0; n < len(table.Fields); n++ {
			field := table.Fields[n]
			protoFields = append(protoFields, render.ProtoField{
				Number:  n + 1,
				Name:    field.Name,
				Comment: field.Comment,
				Type:    field.Type,
			})
		}

		if err := WriteServiceProto(
			outputPath,
			serviceType,
			targetModuleName, sourceModuleName, moduleVersion,
			table.Name, table.Comment,
			protoFields,
		); err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

func WriteDataPackageCode(
	outputPath string,
	orm string,
	projectName string,
	serviceName string,
	name string,
	moduleName, moduleVersion string,
	protoFields []ProtoField,
) error {
	var copyProtoFields render.ProtoFieldArray
	for _, field := range protoFields {
		if field.Type == "" {
			continue
		}

		copyProtoField := render.ProtoField{
			Number:  field.Number,
			Name:    field.Name,
			Type:    field.Type,
			Comment: field.Comment,
		}
		copyProtoFields = append(copyProtoFields, copyProtoField)
	}

	data := render.DataTemplateData{
		Project: projectName,

		Service: serviceName,
		Name:    name,

		Module:  moduleName,
		Version: moduleVersion,

		Fields: copyProtoFields,
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
