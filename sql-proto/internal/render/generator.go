package render

import (
	"os"
	"path/filepath"
	"strings"

	"entgo.io/contrib/schemast"
	"github.com/jinzhu/inflection"

	"github.com/tx7do/kratos-cli/sql-proto/internal"
	"github.com/tx7do/kratos-cli/sql-proto/internal/render/templates"
)

func WriteProto(mutations []schemast.Mutator, opts ...internal.ConvertOption) error {
	return nil
}

func writeGrpcServiceProto(outputPath string, data GrpcProtoTemplateData) {
	outputPath = outputPath + "/" + strings.ToLower(data.Module) + "/service/" + strings.ToLower(data.Version)
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + ProtoFilePostfix
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[GrpcProtoTemplateData](outputPath, data, "grpc_proto", string(templates.GrpcProtoTemplateData))
}

func writeRestServiceProto(outputPath string, data RestProtoTemplateData) {
	outputPath = outputPath + "/" + strings.ToLower(data.TargetModule) + "/service/" + strings.ToLower(data.Version)
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + RestProtoFilePrefix + strings.ToLower(data.Name) + ProtoFilePostfix
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[RestProtoTemplateData](outputPath, data, "rest_proto", string(templates.RestProtoTemplateData))
}

func writeEntDataPackageCode(outputPath string, data DataTemplateData) {
	outputPath = outputPath + "/data/"
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[DataTemplateData](outputPath, data, "ent_data", string(templates.EntDataTemplateData))
}

func writeGormDataPackageCode(outputPath string, data DataTemplateData) {
	outputPath = outputPath + "/data/"
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[DataTemplateData](outputPath, data, "gorm_data", string(templates.GormDataTemplateData))
}

func writeGrpcServicePackageCode(outputPath string, data ServiceTemplateData) {
	outputPath = outputPath + "/service/"
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + GoFilePostfix
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[ServiceTemplateData](outputPath, data, "service", string(templates.ServiceTemplateData))
}

func writeInitWireCode(outputPath string, data InitWireTemplateData) {
	outputPath = outputPath + "/" + data.Package + "/"
	outputPath = filepath.Clean(outputPath)
	_ = os.MkdirAll(outputPath, os.ModePerm)

	goFileName := outputPath + "/" + "init" + GoFilePostfix

	_ = renderTemplate[InitWireTemplateData](goFileName, data, "init_"+data.Package, string(templates.InitTemplateData))
}

func writeWireCode(outputPath string, data WireTemplateData) {
	_ = os.MkdirAll(outputPath, os.ModePerm)

	goFileName := outputPath + "/" + "wire" + GoFilePostfix
	goFileName = filepath.Clean(goFileName)

	data.Service = strings.ToLower(data.Service)

	_ = renderTemplate[WireTemplateData](goFileName, data, "wire_"+data.Project, string(templates.WireTemplateData))
}

func writeMainCode(outputPath string, data MainTemplateData) {
	_ = os.MkdirAll(outputPath, os.ModePerm)

	goFileName := outputPath + "/" + "main" + GoFilePostfix
	goFileName = filepath.Clean(goFileName)

	data.Service = snakeToPascal(data.Service)

	_ = renderTemplate[MainTemplateData](goFileName, data, "main_"+data.Project, string(templates.MainTemplateData))
}

func writeServerPackageCode(outputPath string, data ServerTemplateData) {
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
