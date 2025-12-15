package generators

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/tx7do/go-utils/stringcase"
	"github.com/tx7do/kratos-cli/generatos/templates/golang"
)

// GoGenerator 使用 TemplateEngine 渲染并将结果写入磁盘
type GoGenerator struct {
	Engine TemplateEngine
}

// NewGoGenerator 创建生成器，engine 可为 nil（需要在调用前设置）
func NewGoGenerator() *GoGenerator {
	var funcMap = template.FuncMap{
		"newlineIf": func(condition bool) string {
			if condition {
				return "\n"
			}
			return ""
		},
		"newline": func() string { return "\n" },

		"upper": strings.ToUpper, // 转换为大写
		"lower": strings.ToLower, // 转换为小写

		"camel":  stringcase.LowerCamelCase, // 转换为 camelCase
		"pascal": stringcase.ToPascalCase,   // 转换为 PascalCase
		"kebab":  stringcase.KebabCase,      // 转换为 kebab-case
		"snake":  stringcase.SnakeCase,      // 转换为 snake_case
	}

	templateEngine, _ := NewEmbeddedTemplateEngineFromMap(golang.TemplateMap, funcMap)

	g := &GoGenerator{
		Engine: templateEngine,
	}

	return g
}

// NewGoGeneratorWithEngine 使用指定的引擎创建生成器
func NewGoGeneratorWithEngine(engine TemplateEngine) *GoGenerator {
	g := &GoGenerator{
		Engine: engine,
	}
	return g
}

// Generate 渲染 tplName 并写入 opts.OutDir 下。
// 规则：如果 tplName 以 .tpl 或 .tmpl 结尾，会在输出文件名中去掉该后缀。
func (g *GoGenerator) Generate(opts Options, tplName string) error {
	if g.Engine == nil {
		return os.ErrInvalid
	}

	// 合并数据：以 opts.Vars 为基础，注入常用字段
	data := map[string]any{}
	if opts.Vars != nil {
		for k, v := range opts.Vars {
			data[k] = v
		}
	}
	// 常用上下文
	data["Module"] = opts.Module
	data["ProjectName"] = opts.ProjectName
	data["OutDir"] = opts.OutDir

	// 渲染
	outBytes, err := g.Engine.Render(tplName, data)
	if err != nil {
		return err
	}

	// 计算默认输出名称（保持相对目录并去掉模板后缀）
	defaultOutName := tplName
	if strings.HasSuffix(defaultOutName, ".tpl") {
		defaultOutName = strings.TrimSuffix(defaultOutName, ".tpl")
	} else if strings.HasSuffix(defaultOutName, ".tmpl") {
		defaultOutName = strings.TrimSuffix(defaultOutName, ".tmpl")
	}
	defaultOutName = filepath.FromSlash(defaultOutName)

	// 如果用户指定 OutputName，优先处理
	finalRel := defaultOutName
	if opts.OutputName != "" {
		user := filepath.FromSlash(opts.OutputName)
		// 判断 user 是否包含目录部分
		if filepath.Dir(user) == "." || user == "." {
			// 仅文件名：保留模板的目录结构，替换基础名
			baseDir := filepath.Dir(defaultOutName)
			finalRel = filepath.Join(baseDir, user)
		} else {
			// 含目录：直接使用用户提供的相对路径
			finalRel = user
		}
	}

	outPath := filepath.Join(opts.OutDir, finalRel)
	if err = os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, outBytes, 0o644)
}
