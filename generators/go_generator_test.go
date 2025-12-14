package generators

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoGenerator_GenerateCreatesFile(t *testing.T) {
	tmp := t.TempDir()

	engine, err := NewEmbeddedTemplateEngineFromMap(map[string][]byte{
		"pkg/main.tpl": []byte("package {{.Module}}\n// {{.ProjectName}}\nconst V = {{.Value}}"),
	}, nil)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	g := NewGoGenerator(engine)
	opts := Options{
		Module:      "github.com/example/mod",
		ProjectName: "demo",
		OutDir:      tmp,
		Vars: map[string]interface{}{
			"Value": 42,
		},
	}

	if err := g.Generate(opts, "pkg/main.tpl"); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	outPath := filepath.Join(tmp, "pkg", "main")
	b, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file failed: %v", err)
	}
	got := string(b)
	if !strings.Contains(got, "package github.com/example/mod") || !strings.Contains(got, "// demo") || !strings.Contains(got, "V = 42") {
		t.Fatalf("output content unexpected: %q", got)
	}
}

func TestGoGenerator_RemovesSuffixAndCreatesNestedDirs(t *testing.T) {
	tmp := t.TempDir()

	engine, err := NewEmbeddedTemplateEngineFromMap(map[string][]byte{
		"a/b/c/file.tmpl": []byte("name: {{.ProjectName}}\n"),
	}, nil)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	g := NewGoGenerator(engine)
	opts := Options{
		ProjectName: "nested",
		OutDir:      tmp,
	}

	if err := g.Generate(opts, "a/b/c/file.tmpl"); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	outPath := filepath.Join(tmp, "a", "b", "c", "file")
	b, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file failed: %v", err)
	}
	if strings.TrimSpace(string(b)) != "name: nested" {
		t.Fatalf("unexpected file content: %q", string(b))
	}
}

func TestGoGenerator_NoEngineReturnsError(t *testing.T) {
	g := NewGoGenerator(nil)
	opts := Options{OutDir: t.TempDir()}

	err := g.Generate(opts, "doesnotmatter.tpl")
	if err == nil {
		t.Fatalf("expected error when engine is nil, got nil")
	}
	if !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("expected os.ErrInvalid, got: %v", err)
	}
}

func TestGoGenerator_OverwriteExistingFile(t *testing.T) {
	tmp := t.TempDir()

	engine, err := NewEmbeddedTemplateEngineFromMap(map[string][]byte{
		"dup.tpl": []byte("value: {{.Value}}"),
	}, nil)
	if err != nil {
		t.Fatalf("failed to create engine: %v", err)
	}

	// create existing file with different content
	outPath := filepath.Join(tmp, "dup")
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(outPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write initial file failed: %v", err)
	}

	g := NewGoGenerator(engine)
	opts := Options{
		OutDir: tmp,
		Vars: map[string]interface{}{
			"Value": "new",
		},
	}

	if err := g.Generate(opts, "dup.tpl"); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	b, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output file failed: %v", err)
	}
	if strings.TrimSpace(string(b)) != "value: new" {
		t.Fatalf("file was not overwritten as expected, got: %q", string(b))
	}
}
