package sqlproto

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/tx7do/kratos-cli/sql-proto/internal"
	"github.com/tx7do/kratos-cli/sql-proto/internal/mux"
)

// Convert converts the database schema into a protocol buffer definition.
func Convert(
	ctx context.Context,
	dsn, outputPath *string,
	moduleName, sourceModuleName, moduleVersion *string,
	serviceType *string,
	includeTables, excludeTables []string,
) error {
	if outputPath == nil {
		return errors.New("sqlproto: proto file output path is nil")
	}
	if dsn == nil {
		return errors.New("sqlproto: dsn is nil")
	}
	if moduleName == nil {
		return errors.New("sqlproto: proto module is nil")
	}

	_ = os.MkdirAll(*outputPath, os.ModePerm)

	convertDriver, err := mux.Default.OpenConvert(*dsn)
	if err != nil {
		log.Fatalf("sqlproto: failed to create import driver - %v", err)
	}
	defer convertDriver.Close()

	i, err := internal.NewConvert(
		internal.WithIncludedTables(includeTables),
		internal.WithExcludedTables(excludeTables),
		internal.WithDriver(convertDriver),
	)
	if err != nil {
		log.Fatalf("sqlproto: create importer failed: %v", err)
	}

	tableDatas, err := i.SchemaTables(ctx)
	if err != nil {
		log.Fatalf("sqlproto: schema import failed - %v", err)
	}

	var opts []internal.ConvertOption
	if serviceType != nil {
		opts = append(opts, internal.WithServiceType(*serviceType))
	}
	if outputPath != nil {
		opts = append(opts, internal.WithProtoPath(*outputPath))
	}
	if moduleName != nil {
		opts = append(opts, internal.WithModuleName(*moduleName))
	}
	if sourceModuleName != nil {
		opts = append(opts, internal.WithSourceModuleName(*sourceModuleName))
	}
	if moduleVersion != nil {
		opts = append(opts, internal.WithModuleVersion(*moduleVersion))
	}
	if err = internal.WriteProto(tableDatas, opts...); err != nil {
		log.Fatalf("sqlproto: schema writing failed - %v", err)
	}

	return nil
}

// filterTables filters the provided tables based on the include and exclude lists.
func filterTables(tables []*internal.TableData, includeTables, excludeTables []string) []*internal.TableData {
	if len(excludeTables) == 0 && len(includeTables) == 0 {
		return tables
	}

	tableSet := make(map[string]*internal.TableData)
	for _, table := range tables {
		tableSet[table.Name] = table
	}

	filtered := make([]*internal.TableData, 0, len(includeTables))
	for _, tableName := range includeTables {
		if _, exists := tableSet[tableName]; exists {
			excluded := false
			for _, excludedTable := range excludeTables {
				if tableName == excludedTable {
					excluded = true
					break
				}
			}
			if !excluded {
				filtered = append(filtered, tableSet[tableName])
			}
		}
	}

	return filtered
}
