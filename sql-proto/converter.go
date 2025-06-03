package sqlproto

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/tx7do/kratos-cli/sql-proto/internal"
	"github.com/tx7do/kratos-cli/sql-proto/internal/mux"
)

type TableDataArray []*internal.TableData

// Convert converts the database schema into a protocol buffer definition.
func Convert(
	ctx context.Context,
	dsn, outputPath *string,
	moduleName, sourceModuleName, moduleVersion *string,
	serviceType *string,
	includeTables, excludeTables []string,
	exportProto bool,
) (TableDataArray, error) {
	if outputPath == nil {
		return nil, errors.New("sqlproto: proto file output path is nil")
	}
	if dsn == nil {
		return nil, errors.New("sqlproto: dsn is nil")
	}
	if moduleName == nil {
		return nil, errors.New("sqlproto: proto module is nil")
	}

	_ = os.MkdirAll(*outputPath, os.ModePerm)

	convertDriver, err := mux.Default.OpenConvert(*dsn)
	if err != nil {
		log.Fatalf("sqlproto: failed to create import driver - %v", err)
		return nil, err
	}
	defer convertDriver.Close()

	i, err := internal.NewConvert(
		internal.WithIncludedTables(includeTables),
		internal.WithExcludedTables(excludeTables),
		internal.WithDriver(convertDriver),
	)
	if err != nil {
		log.Fatalf("sqlproto: create importer failed: %v", err)
		return nil, err
	}

	tableDatas, err := i.SchemaTables(ctx)
	if err != nil {
		log.Fatalf("sqlproto: schema import failed - %v", err)
		return nil, err
	}

	if exportProto {
		if err = WriteServicesProto(
			*outputPath,
			*serviceType,
			*moduleName, *sourceModuleName, *moduleVersion,
			tableDatas,
		); err != nil {
			log.Fatalf("sqlproto: schema writing failed - %v", err)
		}
	}

	return tableDatas, nil
}
