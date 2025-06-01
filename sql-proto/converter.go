package sqlproto

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/tx7do/kratos-cli/sql-proto/internal"
	"github.com/tx7do/kratos-cli/sql-proto/internal/mux"
	"github.com/tx7do/kratos-cli/sql-proto/internal/render"
)

// Convert converts the database schema into a protocol buffer definition.
func Convert(ctx context.Context, drv, dsn, outputPath *string, tables, excludeTables []string) error {
	if outputPath == nil {
		return errors.New("sqlproto: proto file output path is nil")
	}
	if dsn == nil {
		return errors.New("sqlproto: dsn is nil")
	}

	_ = os.MkdirAll(*outputPath, os.ModePerm)

	convertDriver, err := mux.Default.OpenConvert(*dsn)
	if err != nil {
		log.Fatalf("sqlproto: failed to create import driver - %v", err)
	}
	defer convertDriver.Close()

	i, err := internal.NewConvert(
		internal.WithTables(tables),
		internal.WithExcludedTables(excludeTables),
		internal.WithDriver(convertDriver),
	)
	if err != nil {
		log.Fatalf("sqlproto: create importer failed: %v", err)
	}

	mutations, err := i.SchemaMutations(ctx)
	if err != nil {
		log.Fatalf("sqlproto: schema import failed - %v", err)
	}

	if err = render.WriteProto(mutations, internal.WithProtoPath(*outputPath)); err != nil {
		log.Fatalf("sqlproto: schema writing failed - %v", err)
	}

	return nil
}
