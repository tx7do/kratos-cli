package render

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tx7do/kratos-cli/sql-proto/internal/render/templates"
)

// WriteGrpcServiceProto GRPC service proto file
func WriteGrpcServiceProto(outputPath string, data GrpcProtoTemplateData) {
	outputPath = outputPath + "/" + strings.ToLower(data.Module) + "/service/" + strings.ToLower(data.Version)
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + ProtoFilePostfix
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[GrpcProtoTemplateData](outputPath, data, "grpc_proto", string(templates.GrpcProtoTemplateData))
}

// WriteRestServiceProto REST service proto file
func WriteRestServiceProto(outputPath string, data RestProtoTemplateData) {
	outputPath = outputPath + "/" + strings.ToLower(data.TargetModule) + "/service/" + strings.ToLower(data.Version)
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + RestProtoFilePrefix + strings.ToLower(data.Name) + ProtoFilePostfix
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[RestProtoTemplateData](outputPath, data, "rest_proto", string(templates.RestProtoTemplateData))
}

// WriteEntDataPackageCode writes the Ent data package code to the specified output path.
func WriteEntDataPackageCode(outputPath string, data DataTemplateData) {
	outputPath = outputPath + "/data/"
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[DataTemplateData](outputPath, data, "ent_data", string(templates.EntDataTemplateData))
}

// WriteGormDataPackageCode writes the GORM data package code to the specified output path.
func WriteGormDataPackageCode(outputPath string, data DataTemplateData) {
	outputPath = outputPath + "/data/"
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[DataTemplateData](outputPath, data, "gorm_data", string(templates.GormDataTemplateData))
}

// WriteGrpcServicePackageCode writes the GRPC service package code to the specified output path.
func WriteGrpcServicePackageCode(outputPath string, data ServiceTemplateData) {
	outputPath = outputPath + "/service/"
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[ServiceTemplateData](outputPath, data, "service", string(templates.ServiceTemplateData))
}

// WriteInitWireCode writes the initialization wire code to the specified output path.
func WriteInitWireCode(outputPath string, data InitWireTemplateData) {
	outputPath = outputPath + "/" + data.Package + "/"
	outputPath = filepath.Clean(outputPath)
	_ = os.MkdirAll(outputPath, os.ModePerm)

	goFileName := outputPath + "/" + "init" + GoFilePostfix

	_ = renderTemplate[InitWireTemplateData](goFileName, data, "init_"+data.Package, string(templates.InitTemplateData))
}

// WriteWireCode writes the wire code to the specified output path.
func WriteWireCode(outputPath string, data WireTemplateData) {
	_ = os.MkdirAll(outputPath, os.ModePerm)

	goFileName := outputPath + "/" + "wire" + GoFilePostfix
	goFileName = filepath.Clean(goFileName)

	data.Service = strings.ToLower(data.Service)

	_ = renderTemplate[WireTemplateData](goFileName, data, "wire_"+data.Project, string(templates.WireTemplateData))
}

// WriteMainCode writes the main code to the specified output path.
func WriteMainCode(outputPath string, data MainTemplateData) {
	_ = os.MkdirAll(outputPath, os.ModePerm)

	goFileName := outputPath + "/" + "main" + GoFilePostfix
	goFileName = filepath.Clean(goFileName)

	data.Service = snakeToPascal(data.Service)

	_ = renderTemplate[MainTemplateData](goFileName, data, "main_"+data.Project, string(templates.MainTemplateData))
}

// WriteServerPackageCode writes the server package code to the specified output path.
func WriteServerPackageCode(outputPath string, data ServerTemplateData) {
	data.Service = strings.ToLower(data.Service)

	outputPath = outputPath + "/server/"
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	switch data.Type {
	case "grpc":
		outputPath = outputPath + "/" + "grpc" + GoFilePostfix
		outputPath = filepath.Clean(outputPath)
		_ = renderTemplate[ServerTemplateData](outputPath, data, "grpc_server_"+data.Project, string(templates.GrpcTemplateServerData))

	case "rest":
		outputPath = outputPath + "/" + "rest" + GoFilePostfix
		outputPath = filepath.Clean(outputPath)
		_ = renderTemplate[ServerTemplateData](outputPath, data, "rest_server_"+data.Project, string(templates.RestTemplateServerData))
	}
}
