package render

import (
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var newlineIfFuncMap = template.FuncMap{
	"newlineIf": func(condition bool) string {
		if condition {
			return "\n"
		}
		return ""
	},
}

var newlineFuncMap = template.FuncMap{
	"newline": func() string { return "\n" },
}

// renderTemplate renders a Protobuf template to a file at the specified output path.
func renderTemplate[T any](outputFileName string, data T, templateName, templateData string) error {
	f, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
		return err
	}
	defer f.Close()

	tmpl := template.
		Must(
			template.New(templateName).
				Funcs(newlineIfFuncMap).
				Funcs(newlineFuncMap).
				Parse(templateData),
		)

	if err = tmpl.Execute(f, data); err != nil {
		log.Fatalf("failed to execute template: %v", err)
		return err
	}

	log.Println(outputFileName, "generated")
	return nil
}

// snakeToCamel converts a snake_case string to camelCase.
func snakeToCamel(snake string) string {
	parts := strings.Split(snake, "_")
	titleCaser := cases.Title(language.Und)
	for i := 1; i < len(parts); i++ {
		parts[i] = titleCaser.String(parts[i])
	}
	return strings.Join(parts, "")
}

// snakeToPascal converts a snake_case string to PascalCase.
func snakeToPascal(snake string) string {
	parts := strings.Split(snake, "_")
	titleCaser := cases.Title(language.Und)
	for i := 0; i < len(parts); i++ { // 从第一个部分开始转换
		parts[i] = titleCaser.String(parts[i])
	}
	return strings.Join(parts, "")
}

func snakeToPascalPlus(snake string) string {
	parts := strings.Split(snake, "_")
	titleCaser := cases.Title(language.Und)
	for i := 0; i < len(parts); i++ {
		part := titleCaser.String(parts[i])
		if strings.ToLower(part) == "id" {
			part = "ID"
		}
		parts[i] = part
	}
	return strings.Join(parts, "")
}

func snakeToKebab(snake string) string {
	return strings.ReplaceAll(snake, "_", "-")
}

// camelToSnake 将字符串从 CamelCase 转换为 snake_case
func camelToSnake(camel string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(camel, "${1}_${2}")
	return strings.ToLower(snake)
}

func makeEntSetNillableFunc(fieldName string) string {
	inputVar := "req.Data." + snakeToPascal(fieldName)
	return "SetNillable" + snakeToPascalPlus(fieldName) + "(" + inputVar + ")"
}

func makeEntSetNillableFuncWithTransfer(fieldName string, transFunc string) string {
	inputVar := "req.Data." + snakeToPascal(fieldName)
	return "SetNillable" + snakeToPascalPlus(fieldName) + "(" + transFunc + "(" + inputVar + "))"
}
