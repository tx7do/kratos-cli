package entimport

import (
	"context"
	"log"

	"github.com/tx7do/kratos-cli/sql-importer/internal/mux"
)

// Importer imports the schema from the database specified by the DSN and writes it to the schemaPath.
func Importer(ctx context.Context, dsn, schemaPath *string, tables, excludeTables []string) error {
	drv, err := mux.Default.OpenImport(*dsn)
	if err != nil {
		log.Fatalf("entimport: failed to create import driver - %v", err)
	}

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
