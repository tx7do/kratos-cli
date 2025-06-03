package sqlproto

import (
	"strings"

	"github.com/tx7do/kratos-cli/sql-proto/internal/render"
)

type ProtoField render.ProtoField
type ProtoFieldArray render.ProtoFieldArray

func WriteServiceProto(
	outputPath string,
	serviceType string,
	name string,
	comment string,
	targetModuleName, sourceModuleName, moduleVersion string,
	protoFields ProtoFieldArray,
) {
	switch strings.TrimSpace(strings.ToLower(serviceType)) {
	case "grpc":
		data := render.GrpcProtoTemplateData{
			Module:  targetModuleName,
			Name:    name,
			Comment: comment,
			Version: moduleVersion,

			Fields: render.ProtoFieldArray(protoFields),
		}
		render.WriteGrpcServiceProto(outputPath, data)
	case "rest":
		data := render.RestProtoTemplateData{
			SourceModule: sourceModuleName,
			TargetModule: targetModuleName,
			Name:         name,
			Version:      moduleVersion,
			Comment:      comment,
		}
		render.WriteRestServiceProto(outputPath, data)
	}
}

func WriteDataPackageCode(
	outputPath string,
	orm string,
	projectName string,
	serviceName string,
	name string,
	moduleName, moduleVersion string,
	protoFields ProtoFieldArray,
) {
	data := render.DataTemplateData{
		Project: projectName,

		Service: serviceName,
		Name:    name,

		Module:  moduleName,
		Version: moduleVersion,

		Fields: render.ProtoFieldArray(protoFields),
	}
	switch strings.TrimSpace(strings.ToLower(orm)) {
	case "ent":
		render.WriteEntDataPackageCode(outputPath, data)

	case "gorm":
		render.WriteGormDataPackageCode(outputPath, data)
	}
}

func WriteServicePackageCode(
	outputPath string,
	projectName string,
	serviceName string,
	name string,
	targetModuleName, sourceModuleName, moduleVersion string,
	userRepo, isGrpcService bool,
) {
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
	render.WriteServicePackageCode(outputPath, data)
}

func WriteServerPackageCode(
	outputPath string,
	projectName string,
	serviceType string,
	serviceName string,
	services map[string]string,
) {
	data := render.ServerTemplateData{
		Project:  projectName,
		Type:     serviceType,
		Service:  serviceName,
		Services: services,
	}
	render.WriteServerPackageCode(outputPath, data)
}

func WriteInitWireCode(
	outputPath string,

	projectName string,
	postfix string,
	services []string,
) {
	data := render.InitWireTemplateData{
		Package:      projectName,
		Postfix:      postfix,
		ServiceNames: services,
	}
	render.WriteInitWireCode(outputPath, data)
}

func WriteWireCode(
	outputPath string,

	projectName string,
	serviceName string,
) {
	data := render.WireTemplateData{
		Project: projectName,
		Service: serviceName,
	}
	render.WriteWireCode(outputPath, data)
}

func WriteMainCode(
	outputPath string,

	projectName string,
	serviceName string,

	enableREST bool,
	enableGRPC bool,
	enableAsynq bool,
	enableSSE bool,
	enableKafka bool,
	enableMQTT bool,
) {
	data := render.MainTemplateData{
		Project: projectName,
		Service: serviceName,

		EnableREST:  enableREST,
		EnableGRPC:  enableGRPC,
		EnableAsynq: enableAsynq,
		EnableSSE:   enableSSE,
		EnableKafka: enableKafka,
		EnableMQTT:  enableMQTT,
	}
	render.WriteMainCode(outputPath, data)
}
