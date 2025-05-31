package sqlorm

import (
	"context"
	"errors"
	"strings"

	"github.com/tx7do/kratos-cli/sql-orm/internal/ent/entimport"
	"github.com/tx7do/kratos-cli/sql-orm/internal/gorm"
)

func Importer(ctx context.Context, orm string, drv, dsn, schemaPath, daoPath *string, tables, excludeTables []string) error {
	switch strings.ToLower(strings.TrimSpace(orm)) {
	case "ent":
		return entimport.Importer(ctx, dsn, schemaPath, tables, excludeTables)

	case "gorm":
		return gorm.Importer(ctx, drv, dsn, schemaPath, daoPath, tables, excludeTables)

	default:
		return errors.New("sql2orm: unsupported orm type: " + orm)
	}
}
