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

type ProtoFieldArray []ProtoField

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

type GrpcProtoTemplateData struct {
	Name    string // Proto文件名
	Comment string // Proto文件注释
	Version string // Proto版本号

	Module string // 模块名

	Fields ProtoFieldArray // 字段列表
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
