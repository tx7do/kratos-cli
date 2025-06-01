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

const (
	ProtoFilePostfix    = ".proto"
	RestProtoFilePrefix = "i_"
	GoFilePostfix       = ".go"
)

type ProtoField struct {
	Name    string
	Type    string
	Number  int
	Comment string
}

func (f ProtoField) CamelName() string {
	return snakeToCamel(f.Name)
}

func (f ProtoField) SnakeName() string {
	return camelToSnake(f.Name)
}

func (f ProtoField) PascalName() string {
	return snakeToPascal(f.Name)
}

func (f ProtoField) EntPascalName() string {
	return snakeToPascalPlus(f.Name)
}

func (f ProtoField) EntSetNillableFunc() string {
	return makeEntSetNillableFunc(f.Name)
}

type GrpcProtoTemplateData struct {
	Name    string
	Comment string
	Version string

	Module string

	Fields []ProtoField
}

func (d GrpcProtoTemplateData) PascalName() string {
	return snakeToPascal(d.Name)
}
func (d GrpcProtoTemplateData) SnakeName() string {
	return camelToSnake(d.Name)
}

func (d GrpcProtoTemplateData) Package() string {
	return strings.ToLower(d.Module) + ".service." + d.Version
}

type RestProtoTemplateData struct {
	Name    string
	Comment string
	Version string

	SourceModule string
	TargetModule string
}

func (d RestProtoTemplateData) PascalName() string {
	return snakeToPascal(d.Name)
}

func (d RestProtoTemplateData) Path() string {
	// 比如：/admin/v1/users
	return "/" + strings.ToLower(d.TargetModule) + "/" + d.Version + "/" + snakeToKebab(inflection.Plural(d.Name))
}

func (d RestProtoTemplateData) SourceProto() string {
	// 比如：user/service/v1/user.proto
	return strings.ToLower(d.SourceModule) + "/service/" + d.Version + "/" + strings.ToLower(d.Name) + ".proto"
}

func (d RestProtoTemplateData) SourcePackage() string {
	// 比如：user.service.v1
	return strings.ToLower(d.SourceModule) + ".service." + d.Version
}

func (d RestProtoTemplateData) TargetPackage() string {
	return strings.ToLower(d.TargetModule) + ".service." + d.Version
}

type DataTemplateData struct {
	Project string
	Service string
	Name    string

	Module  string
	Version string

	Fields []ProtoField

	UseTimestamp bool
}

func (d DataTemplateData) LowerName() string {
	return strings.ToLower(d.Name)
}

func (d DataTemplateData) PascalName() string {
	return snakeToPascal(d.Name)
}

func (d DataTemplateData) ApiPackage() string {
	return strings.ToLower(d.Module) + strings.ToUpper(d.Version)
}

func (d DataTemplateData) ClassName() string {
	return d.PascalName() + "Repo"
}

type ServiceTemplateData struct {
	Project string
	Service string
	Name    string

	Version string

	SourceApi string
	TargetApi string

	UseRepo bool // 是否使用数据仓库，否则使用GRPC客户端。
	IsGrpc  bool // 是否是GRPC服务。
}

func (d ServiceTemplateData) LowerName() string {
	return strings.ToLower(d.Name)
}

func (d ServiceTemplateData) PascalName() string {
	return snakeToPascal(d.Name)
}

func (d ServiceTemplateData) SourceApiPackage() string {
	return strings.ToLower(d.SourceApi) + strings.ToUpper(d.Version)
}

func (d ServiceTemplateData) TargetApiPackage() string {
	return strings.ToLower(d.TargetApi) + strings.ToUpper(d.Version)
}

func (d ServiceTemplateData) IsSameApi() bool {
	return strings.ToLower(d.SourceApi) == strings.ToLower(d.TargetApi)
}

func (d ServiceTemplateData) ServiceInterface() string {
	if d.IsGrpc {
		return d.TargetApiPackage() + "." + "Unimplemented" + d.PascalName() + "ServiceServer"
	} else {
		return d.TargetApiPackage() + "." + d.LowerName() + "ServiceHTTPServer"
	}
}

func (d ServiceTemplateData) DataSourceType() string {
	if d.UseRepo {
		return "*data." + d.PascalName() + "Repo"
	} else {
		return d.SourceApiPackage() + "." + d.PascalName() + "ServiceClient"
	}
}

func (d ServiceTemplateData) DataSourceVar() string {
	if d.UseRepo {
		return "repo"
	} else {
		return "uc"
	}
}

type InitWireTemplateData struct {
	Package      string
	Postfix      string
	ServiceNames []string
}

type WireTemplateData struct {
	Project string
	Service string
}

type MainTemplateData struct {
	Project string
	Service string

	EnableREST  bool
	EnableGRPC  bool
	EnableAsynq bool
	EnableSSE   bool
	EnableKafka bool
	EnableMQTT  bool
}

type ProtoItem struct {
	Module string
	Name   string
}

type ServerTemplateData struct {
	Project  string
	Service  string
	Type     string
	Services map[string]string
}

func (d ServerTemplateData) Modules() []string {
	modulesMap := make(map[string]struct{})
	for _, module := range d.Services {
		modulesMap[module] = struct{}{}
	}

	var modules []string
	for module := range modulesMap {
		modules = append(modules, module)
	}
	return modules
}

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

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + ".go"
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[DataTemplateData](outputPath, data, "ent_data", string(templates.EntDataTemplateData))
}

func writeGormDataPackageCode(outputPath string, data DataTemplateData) {
	outputPath = outputPath + "/data/"
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + ".go"
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[DataTemplateData](outputPath, data, "gorm_data", string(templates.GormDataTemplateData))
}

func writeGrpcServicePackageCode(outputPath string, data ServiceTemplateData) {
	outputPath = outputPath + "/service/"
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	outputPath = outputPath + "/" + strings.ToLower(data.Name) + ".go"
	outputPath = filepath.Clean(outputPath)

	_ = renderTemplate[ServiceTemplateData](outputPath, data, "service", string(templates.ServiceTemplateData))
}

func writeInitWireCode(outputPath string, data InitWireTemplateData) {
	outputPath = outputPath + "/" + data.Package + "/"
	outputPath = filepath.Clean(outputPath)
	_ = os.MkdirAll(outputPath, os.ModePerm)

	goFileName := outputPath + "/" + "init.go"

	_ = renderTemplate[InitWireTemplateData](goFileName, data, "init_"+data.Package, string(templates.InitTemplateData))
}

func writeWireCode(outputPath string, data WireTemplateData) {
	_ = os.MkdirAll(outputPath, os.ModePerm)

	goFileName := outputPath + "/" + "wire.go"
	goFileName = filepath.Clean(goFileName)

	data.Service = strings.ToLower(data.Service)

	_ = renderTemplate[WireTemplateData](goFileName, data, "wire_"+data.Project, string(templates.WireTemplateData))
}

func writeMainCode(outputPath string, data MainTemplateData) {
	_ = os.MkdirAll(outputPath, os.ModePerm)

	goFileName := outputPath + "/" + "main.go"
	goFileName = filepath.Clean(goFileName)

	data.Service = snakeToPascal(data.Service)

	_ = renderTemplate[MainTemplateData](goFileName, data, "main_"+data.Project, string(templates.MainTemplateData))
}

func writeServerPackageCode(outputPath string, data ServerTemplateData) {
	outputPath = outputPath + "/server/"
	outputPath = filepath.Clean(outputPath)

	_ = os.MkdirAll(outputPath, os.ModePerm)

	goFileName := outputPath + "/" + "grpc.go"
	goFileName = filepath.Clean(goFileName)

	data.Service = strings.ToLower(data.Service)

	_ = renderTemplate[ServerTemplateData](goFileName, data, "server_"+data.Project, string(templates.GrpcTemplateServerData))
}
