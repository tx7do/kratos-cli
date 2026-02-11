package project

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/tx7do/kratos-cli/gowind/internal/pkg"
)

const (
	TemplateModuleName = "github.com/tx7do/go-wind-admin-template"
)

type Project struct {
	Name   string
	Path   string
	Module string
}

// New a project from remote repo.
func (p *Project) New(ctx context.Context, dir string, layout string, branch string) error {
	to := filepath.Join(dir, p.Name)

	repo := pkg.NewRepo(layout, branch)
	if err := repo.CopyTo(ctx, to, p.Name, []string{".git", ".github"}); err != nil {
		return err
	}

	updateCount, err := pkg.ReplaceTemplateInCurrentDir(dir, TemplateModuleName, p.Module)
	if err != nil {
		return err
	}
	log.Printf("Updated %d files.\n", updateCount)

	return nil
}

var repoAddIgnores = []string{
	".git", ".github", "api", "README.md", "LICENSE", "go.mod", "go.sum", "third_party", "openapi.yaml", ".gitignore",
}

func (p *Project) Add(ctx context.Context, dir string, layout string, branch string, mod string, pkgPath string) error {
	to := filepath.Join(dir, p.Name)

	log.Printf("🚀 Add service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)

	pkgPath = fmt.Sprintf("%s/%s", mod, pkgPath)
	repo := pkg.NewRepo(layout, branch)
	err := repo.CopyToV2(ctx, to, pkgPath, repoAddIgnores, []string{filepath.Join(p.Path, "api"), "api"})
	if err != nil {
		return err
	}

	updateCount, err := pkg.ReplaceTemplateInCurrentDir(dir, TemplateModuleName, p.Module)
	if err != nil {
		return err
	}
	log.Printf("Updated %d files.\n", updateCount)

	return nil
}
