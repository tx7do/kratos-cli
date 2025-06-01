package render

import (
	"strings"

	"github.com/jinzhu/inflection"
)

// ProtoField protobuf字段定义
type ProtoField struct {
	Number  int    // 字段编号
	Name    string // 字段名
	Type    string // 字段类型
	Comment string // 字段注释
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

func (d ServiceTemplateData) ClassName() string {
	return d.PascalName() + "Service"
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
