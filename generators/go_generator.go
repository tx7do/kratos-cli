package generators

import (
	"context"
	"fmt"
	"os"

	"github.com/tx7do/go-utils/code_generator"
	"github.com/tx7do/go-utils/stringcase"
	"github.com/tx7do/kratos-cli/generators/templates/golang"
)

// GoGenerator 使用 TemplateEngine 渲染并将结果写入磁盘
type GoGenerator struct {
	*code_generator.CodeGenerator
}

// NewGoGenerator 创建生成器，engine 可为 nil（需要在调用前设置）
func NewGoGenerator() *GoGenerator {
	templateEngine, _ := code_generator.NewEmbeddedTemplateEngineFromMap(golang.TemplateMap, funcMap)

	codeGenerator := code_generator.NewCodeGeneratorWithEngine(templateEngine)

	g := &GoGenerator{
		CodeGenerator: codeGenerator,
	}

	return g
}

func (g *GoGenerator) GenerateMain(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}
	return g.Generate(ctx, opts, "main.tpl")
}

func (g *GoGenerator) GenerateWire(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}
	return g.Generate(ctx, opts, "wire.tpl")
}

func (g *GoGenerator) GenerateInit(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}
	return g.Generate(ctx, opts, "init.tpl")
}

func (g *GoGenerator) GenerateData(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}
	return g.Generate(ctx, opts, "data.tpl")
}

func (g *GoGenerator) GenerateEntClient(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}
	return g.Generate(ctx, opts, "ent_client.tpl")
}

func (g *GoGenerator) GenerateEntRepo(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}

	var modelName string
	if v, ok := opts.Vars["Model"]; ok {
		modelName, _ = v.(string)
	}

	if _, ok := opts.Vars["ClassName"]; !ok {
		opts.Vars["ClassName"] = stringcase.ToPascalCase(modelName) + "Repo"
	}

	if _, ok := opts.Vars["ApiPackageVersion"]; !ok {
		opts.Vars["ApiPackageVersion"] = "v1"
	}

	if opts.OutputName == "" {
		opts.OutputName = stringcase.ToSnakeCase(modelName) + "_repo.go"
	}

	return g.Generate(ctx, opts, "ent_repo.tpl")
}

func (g *GoGenerator) GenerateGormClient(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}
	return g.Generate(ctx, opts, "gorm_client.tpl")
}

func (g *GoGenerator) GenerateGormInit(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}
	return g.Generate(ctx, opts, "gorm_init.tpl")
}

func (g *GoGenerator) GenerateGormRepo(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}

	var modelName string
	if v, ok := opts.Vars["Model"]; ok {
		modelName, _ = v.(string)
	}

	if _, ok := opts.Vars["ClassName"]; !ok {
		opts.Vars["ClassName"] = stringcase.ToPascalCase(modelName) + "Repo"
	}

	if _, ok := opts.Vars["ApiPackageVersion"]; !ok {
		opts.Vars["ApiPackageVersion"] = "v1"
	}

	if opts.OutputName == "" {
		opts.OutputName = stringcase.ToSnakeCase(modelName) + "_repo.go"
	}

	return g.Generate(ctx, opts, "gorm_repo.tpl")
}

func (g *GoGenerator) GenerateGrpcServiceProto(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}

	var modelName string
	if v, ok := opts.Vars["Model"]; ok {
		modelName, _ = v.(string)
	}

	if opts.OutputName == "" {
		opts.OutputName = stringcase.ToSnakeCase(modelName) + ".proto"
	}

	return g.Generate(ctx, opts, "grpc_proto.tpl")
}

func (g *GoGenerator) GenerateGrpcServer(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}

	return g.Generate(ctx, opts, "grpc_server.tpl")
}

func (g *GoGenerator) GenerateRedisClient(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}

	return g.Generate(ctx, opts, "redis_client.tpl")
}

func (g *GoGenerator) GenerateRestServiceProto(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}

	var modelName string
	if v, ok := opts.Vars["Model"]; ok {
		modelName, _ = v.(string)
	}

	if opts.OutputName == "" {
		opts.OutputName = "i_" + stringcase.ToSnakeCase(modelName) + ".proto"
	}

	return g.Generate(ctx, opts, "rest_proto.tpl")
}

func (g *GoGenerator) GenerateRestServer(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}

	return g.Generate(ctx, opts, "rest_server.tpl")
}

func (g *GoGenerator) GenerateService(ctx context.Context, opts code_generator.Options) (outputPath string, err error) {
	if g.CodeGenerator == nil {
		return "", os.ErrInvalid
	}

	var modelName string
	if v, ok := opts.Vars["Model"]; ok {
		modelName, _ = v.(string)
	}

	var isGrpcService bool
	if v, ok := opts.Vars["IsGrpc"]; ok {
		isGrpcService, _ = v.(bool)
	}

	if _, ok := opts.Vars["ClassName"]; !ok {

		opts.Vars["ClassName"] = stringcase.ToPascalCase(modelName) + "Service"
	}

	if _, ok := opts.Vars["TargetApiPackageVersion"]; !ok {
		opts.Vars["TargetApiPackageVersion"] = "v1"
	}
	if _, ok := opts.Vars["SourceApiPackageVersion"]; !ok {
		opts.Vars["SourceApiPackageVersion"] = "v1"
	}

	if _, ok := opts.Vars["SourceApiPackage"]; !ok {
		opts.Vars["SourceApiPackage"] = stringcase.LowerCamelCase(opts.Vars["SourceApiPackageName"].(string)) + stringcase.UpperCamelCase(opts.Vars["SourceApiPackageVersion"].(string))
	}
	if _, ok := opts.Vars["TargetApiPackage"]; !ok {
		opts.Vars["TargetApiPackage"] = stringcase.LowerCamelCase(opts.Vars["TargetApiPackageName"].(string)) + stringcase.UpperCamelCase(opts.Vars["TargetApiPackageVersion"].(string))
	}

	if _, ok := opts.Vars["ServiceInterface"]; !ok {
		if isGrpcService {
			opts.Vars["ServiceInterface"] = fmt.Sprintf("%s.Unimplemented%sServiceServer",
				opts.Vars["TargetApiPackage"].(string),
				stringcase.ToPascalCase(opts.Vars["Service"].(string)))
		} else {
			opts.Vars["ServiceInterface"] = fmt.Sprintf("%s.%sServiceHTTPServer",
				opts.Vars["TargetApiPackage"].(string),
				stringcase.ToPascalCase(opts.Vars["Model"].(string)))
		}
	}

	if _, ok := opts.Vars["DataSourceVar"]; !ok {
		if isGrpcService {
			opts.Vars["DataSourceVar"] = stringcase.LowerCamelCase(opts.Vars["Model"].(string)) + "Repo"
		} else {
			opts.Vars["DataSourceVar"] = stringcase.LowerCamelCase(opts.Vars["Model"].(string)) + "ServiceClient"
		}
	}
	if _, ok := opts.Vars["DataSourceType"]; !ok {
		if isGrpcService {
			opts.Vars["DataSourceType"] = "*data." + stringcase.UpperCamelCase(opts.Vars["Model"].(string)) + "Repo"
		} else {
			opts.Vars["DataSourceType"] = fmt.Sprintf("%s.%sServiceClient",
				opts.Vars["SourceApiPackage"].(string),
				stringcase.UpperCamelCase(opts.Vars["Model"].(string)))
		}
	}

	if _, ok := opts.Vars["IsSameApi"]; !ok {
		opts.Vars["IsSameApi"] = opts.Vars["SourceApiPackage"].(string) == opts.Vars["TargetApiPackage"].(string)
	}

	if _, ok := opts.Vars["UseRepo"]; !ok {
		if isGrpcService {
			opts.Vars["UseRepo"] = true
		}
	}

	if opts.OutputName == "" {
		opts.OutputName = stringcase.ToSnakeCase(modelName) + "_service.go"
	}

	return g.Generate(ctx, opts, "service.tpl")
}
