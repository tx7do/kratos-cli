package generators

import (
	"os"
	"path/filepath"
	"strings"
)

// GoGenerator 使用 TemplateEngine 渲染并将结果写入磁盘
type GoGenerator struct {
	Engine TemplateEngine
}

// NewGoGenerator 创建生成器，engine 可为 nil（需要在调用前设置）
func NewGoGenerator(templateEngine TemplateEngine) *GoGenerator {
	return &GoGenerator{
		Engine: templateEngine,
	}
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

	// 计算输出路径：保持 tplName 的相对目录结构，去掉模板后缀
	outName := tplName
	if strings.HasSuffix(outName, ".tpl") {
		outName = strings.TrimSuffix(outName, ".tpl")
	} else if strings.HasSuffix(outName, ".tmpl") {
		outName = strings.TrimSuffix(outName, ".tmpl")
	}
	// 使用 filepath.Join 并保持子目录
	outPath := filepath.Join(opts.OutDir, filepath.FromSlash(outName))
	if err = os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(outPath, outBytes, 0o644)
}
