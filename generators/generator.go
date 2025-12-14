package generators

type Options struct {
	Module      string
	ProjectName string
	OutDir      string
	Vars        map[string]interface{}
}

// Generator 通用生成器：渲染指定模板并写入输出
type Generator interface {
	// Generate 根据选项渲染指定模板并写入输出目录
	Generate(opts Options, tplName string) error
}

// TemplateEngine 模板引擎接口：加载/渲染/列出模板
type TemplateEngine interface {
	// Render 渲染指定模板并返回结果
	Render(tplName string, data any) ([]byte, error)

	// ListTemplates 列出可用的模板名称
	ListTemplates() []string
}
