package sqlkratos

import (
	"context"
	"fmt"

	sqlorm "github.com/tx7do/kratos-cli/sql-orm"
	sqlproto "github.com/tx7do/kratos-cli/sql-proto"

	"github.com/tx7do/kratos-cli/sql-kratos/internal"
)

func Generate(ctx context.Context, opts internal.GeneratorOptions) error {
	var err error

	var tables sqlproto.TableDataArray

	// 生成 Protobuf schema
	if tables, err = generateProtobufCode(ctx, opts); err != nil {
		return err
	}

	services := make([]string, 0)
	servicePackageMap := make(map[string]string)
	for _, table := range tables {
		if len(table.Fields) == 0 {
			continue
		}

		name := table.Name

		services = append(services, name)
		servicePackageMap[name] = opts.ModuleName
	}

	var useGrpc bool
	for _, server := range opts.Servers {
		if server == "grpc" {
			useGrpc = true
			break
		}
	}

	serviceRootPath := opts.OutputPath + "/service"

	// 生成ORM代码
	if opts.GenerateServer {
		if err = generateOrmCode(ctx, opts, serviceRootPath); err != nil {
			return err
		}
	}

	// 生成data层代码
	if opts.GenerateData {
		dataPackagePath := fmt.Sprintf("%s/app/%s/service/internal/", opts.OutputPath, opts.ModuleName)
		if err = generateDataPackageCode(
			dataPackagePath,
			opts.OrmType,
			opts.ProjectName,
			opts.ServiceName,
			opts.ModuleName, opts.ModuleVersion,
			tables,
			services,
		); err != nil {
			return err
		}
	}

	// 生成service层代码
	if opts.GenerateService {
		servicePackagePath := fmt.Sprintf("%s/app/%s/service/internal/", opts.OutputPath, opts.ModuleName)
		if err = generateServicePackageCode(
			servicePackagePath,
			opts.ProjectName,
			opts.ServiceName,
			opts.ModuleName, opts.SourceModuleName, opts.ModuleVersion,
			opts.UseRepo, useGrpc,
			tables,
			services,
		); err != nil {
			return err
		}
	}

	// 生成server层代码
	if opts.GenerateServer {
		serverPackagePath := fmt.Sprintf("%s/app/%s/service/internal/", opts.OutputPath, opts.ModuleName)
		if err = generateServerPackageCode(
			serverPackagePath,
			opts.ProjectName,
			opts.ServiceName,
			servicePackageMap,
			opts.Servers,
		); err != nil {
			return err
		}
	}

	// 生成main包代码
	if opts.GenerateMain {
		mainPackagePath := fmt.Sprintf("%s/app/%s/service/cmd/server", opts.OutputPath, opts.ModuleName)
		if err = generateMainPackageCode(
			mainPackagePath,
			opts.ProjectName,
			opts.ServiceName,
			opts.Servers,
		); err != nil {
			return err
		}
	}

	return nil
}

// generateProtobufCode generates the Protobuf code from the database schema.
func generateProtobufCode(ctx context.Context, opts internal.GeneratorOptions) (sqlproto.TableDataArray, error) {
	var err error
	var tables sqlproto.TableDataArray

	protoPath := opts.OutputPath + "/api/protos/"

	for _, server := range opts.Servers {
		if server != "grpc" && server != "rest" {
			continue
		}

		if tables, err = sqlproto.Convert(
			ctx,
			&opts.Source,
			&opts.OutputPath,
			&protoPath,
			&opts.SourceModuleName,
			&opts.ModuleVersion,
			&server,
			opts.IncludedTables,
			opts.ExcludedTables,
			opts.GenerateProto,
		); err != nil {
			return nil, err
		}
	}

	return tables, nil
}

// generateOrmCode generates the ORM code based on the specified ORM type.
func generateOrmCode(
	ctx context.Context,
	opts internal.GeneratorOptions,
	serviceRootPath string,
) error {
	var err error

	var schemaPath string
	var daoPath string
	switch opts.OrmType {
	case "ent":
		schemaPath = serviceRootPath + "/data/ent/schema"
	case "gorm":
		schemaPath = serviceRootPath + "/data/gorm/schema"
		daoPath = serviceRootPath + "/data/gorm/dao"
	}

	if err = sqlorm.Importer(
		ctx,
		opts.OrmType,
		&opts.Driver,
		&opts.Source,
		&schemaPath,
		&daoPath,
		opts.IncludedTables,
		opts.ExcludedTables,
	); err != nil {
		return err
	}

	return nil
}

func generateServerPackageCode(
	outputPath string,
	projectName string,
	serviceName string,
	servicePackageMap map[string]string,
	servers []string,
) error {

	for _, server := range servers {
		if err := WriteServerPackageCode(
			outputPath,
			projectName, server, serviceName,
			servicePackageMap,
		); err != nil {
			return err
		}
	}

	return WriteInitWireCode(outputPath, projectName, "Server", servers)
}

func generateServicePackageCode(
	outputPath string,
	projectName, serviceName string,
	targetModuleName, sourceModuleName, moduleVersion string,
	userRepo, isGrpcService bool,
	tables sqlproto.TableDataArray,
	services []string,
) error {

	for _, table := range tables {
		if len(table.Fields) == 0 {
			continue
		}

		name := table.Name

		if err := WriteServicePackageCode(
			outputPath,
			projectName, serviceName,
			name,
			targetModuleName, sourceModuleName, moduleVersion,
			userRepo, isGrpcService,
		); err != nil {
			return err
		}
	}

	return WriteInitWireCode(outputPath, projectName, "Service", services)
}

func generateDataPackageCode(
	outputPath string,
	orm string,
	projectName string, serviceName string,
	moduleName, moduleVersion string,
	tables sqlproto.TableDataArray,
	services []string,
) error {
	if len(tables) == 0 {
		return nil
	}

	var dataFields []DataField
	for _, table := range tables {
		if len(table.Fields) == 0 {
			continue
		}

		name := table.Name

		dataFields = make([]DataField, 0)
		for _, field := range table.Fields {
			if field.Type == "" {
				continue
			}

			dataField := DataField{
				Name: field.Name,
				Type: field.Type,
				//Null:    field.Null,
				Comment: field.Comment,
			}
			dataFields = append(dataFields, dataField)
		}

		if err := WriteDataPackageCode(
			outputPath,
			orm,
			projectName, serviceName, name,
			moduleName, moduleVersion,
			dataFields,
		); err != nil {
			return err
		}
	}

	return WriteInitWireCode(outputPath, projectName, "Repo", services)
}

func generateMainPackageCode(
	outputPath string,

	projectName string, serviceName string,

	servers []string,
) error {
	if err := WriteMainCode(
		outputPath,
		projectName, serviceName,
		servers,
	); err != nil {
		return err
	}

	return WriteWireCode(
		outputPath,
		projectName, serviceName,
	)
}
