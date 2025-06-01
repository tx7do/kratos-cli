package entimport

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/tx7do/kratos-cli/sql-orm/internal/ent/mux"
)

// Importer imports the schema from the database specified by the DSN and writes it to the schemaPath.
func Importer(ctx context.Context, dsn, schemaPath *string, tables, excludeTables []string) error {
	if schemaPath == nil {
		return errors.New("entimport: schema path is nil")
	}
	if dsn == nil {
		return errors.New("entimport: dsn is nil")
	}

	_ = os.MkdirAll(*schemaPath, os.ModePerm)

	drv, err := mux.Default.OpenImport(*dsn)
	if err != nil {
		log.Fatalf("entimport: failed to create import driver - %v", err)
	}
	defer drv.Close()

	i, err := NewImport(
		WithTables(tables),
		WithExcludedTables(excludeTables),
		WithDriver(drv),
	)
	if err != nil {
		log.Fatalf("entimport: create importer failed: %v", err)
	}

	mutations, err := i.SchemaMutations(ctx)
	if err != nil {
		log.Fatalf("entimport: schema import failed - %v", err)
	}

	if err = WriteSchema(mutations, WithSchemaPath(*schemaPath)); err != nil {
		log.Fatalf("entimport: schema writing failed - %v", err)
	}

	return nil
}
