package render

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tx7do/kratos-cli/sql-proto/internal/render/templates"
)

// WriteGrpcServiceProto write gRPC service proto file
func WriteGrpcServiceProto(outputPath string, data GrpcProtoTemplateData) error {
	outputPath = outputPath + "/" + strings.ToLower(data.Module) + "/service/" + strings.ToLower(data.Version)
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + ProtoFilePostfix
	outputPath = filepath.Clean(outputPath)

	return renderTemplate[GrpcProtoTemplateData](outputPath, data, "grpc_proto", string(templates.GrpcProtoTemplateData))
}

// WriteRestServiceProto write REST service proto file
func WriteRestServiceProto(outputPath string, data RestProtoTemplateData) error {
	outputPath = outputPath + "/" + strings.ToLower(data.TargetModule) + "/service/" + strings.ToLower(data.Version)
	outputPath = filepath.Clean(outputPath)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}

	outputPath = outputPath + "/" + RestProtoFilePrefix + strings.ToLower(data.Name) + ProtoFilePostfix
	outputPath = filepath.Clean(outputPath)

	return renderTemplate[RestProtoTemplateData](outputPath, data, "rest_proto", string(templates.RestProtoTemplateData))
}
