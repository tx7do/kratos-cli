package sqlkratos

import (
	"context"
	"errors"
	"strings"

	"github.com/tx7do/go-utils/code_generator"
	"github.com/tx7do/go-utils/stringcase"
	"github.com/tx7do/kratos-cli/generators"
)

func WriteDataPackageCode(
	outputPath string,
	orm string,
	projectName string,
	serviceName string,
	name string,
	moduleName, moduleVersion string,
	protoFields []generators.DataField,
) error {
	var copyDataFields generators.DataFieldArray
	for _, field := range protoFields {
		if field.Type == "" {
			continue
		}

		copyDataField := generators.DataField{
			Name:    field.Name,
			Type:    field.Type,
			Comment: field.Comment,
		}
		copyDataFields = append(copyDataFields, copyDataField)
	}

	switch strings.TrimSpace(strings.ToLower(orm)) {
	case "ent":
		if err := writeEntClientCode(outputPath, projectName, serviceName); err != nil {
			return err
		}
		return writeEntRepoCode(outputPath, projectName, serviceName, name, moduleName, moduleVersion, copyDataFields)

	case "gorm":
		if err := writeGormClientCode(outputPath, projectName, serviceName); err != nil {
			return err
		}
		return writeGormRepoCode(outputPath, projectName, serviceName, name, moduleName, moduleVersion, copyDataFields)

	default:
		return errors.New("sqlproto: unsupported orm: " + orm)
	}
}

func writeEntClientCode(
	outputPath string,
	projectModule string,
	serviceName string,
) error {
	g := generators.NewGoGenerator()

	opts := code_generator.Options{
		OutDir: outputPath,
		Module: projectModule,
		Vars: map[string]interface{}{
			"Service": serviceName,
		},
	}

	_, err := g.GenerateEntClient(context.Background(), opts)
	return err
}

func writeGormClientCode(
	outputPath string,
	projectModule string,
	serviceName string,
) error {
	g := generators.NewGoGenerator()

	opts1 := code_generator.Options{
		OutDir: outputPath,
		Module: projectModule,
		Vars: map[string]interface{}{
			"Service": serviceName,
		},
	}
	_, err := g.GenerateGormClient(context.Background(), opts1)
	if err != nil {
		return err
	}

	opts2 := code_generator.Options{
		OutDir: outputPath,
		Module: projectModule,
		Vars: map[string]interface{}{
			"Service": serviceName,
		},
	}
	_, err = g.GenerateGormInit(context.Background(), opts2)
	if err != nil {
		return err
	}

	return nil
}

func writeEntRepoCode(
	outputPath string,
	projectModule string,
	serviceName string,
	model string,
	apiPackageName string,
	apiPackageVersion string,
	fields generators.DataFieldArray,
) error {
	g := generators.NewGoGenerator()

	opts := code_generator.Options{
		OutDir: outputPath,
		Module: projectModule,
		Vars: map[string]interface{}{
			"Service":    serviceName,
			"ApiPackage": stringcase.LowerCamelCase(apiPackageName) + stringcase.UpperCamelCase(apiPackageVersion),
			"Model":      model,
			"Fields":     fields,
		},
	}

	_, err := g.GenerateEntRepo(context.Background(), opts)
	return err
}

func writeGormRepoCode(
	outputPath string,
	projectModule string,
	serviceName string,
	model string,
	apiPackageName string,
	apiPackageVersion string,
	fields generators.DataFieldArray,
) error {
	g := generators.NewGoGenerator()

	opts := code_generator.Options{
		OutDir: outputPath,
		Module: projectModule,
		Vars: map[string]interface{}{
			"Service":    serviceName,
			"ApiPackage": stringcase.LowerCamelCase(apiPackageName) + stringcase.UpperCamelCase(apiPackageVersion),
			"Model":      model,
			"Fields":     fields,
		},
	}

	_, err := g.GenerateGormRepo(context.Background(), opts)
	return err
}

func WriteServicePackageCode(
	outputPath string,
	projectName string,
	serviceName string,
	name string,
	targetModuleName, sourceModuleName, moduleVersion string,
	useRepo, isGrpcService bool,
) error {
	g := generators.NewGoGenerator()

	o := code_generator.Options{
		OutDir: outputPath,
		Module: projectName,
		Vars: map[string]interface{}{
			"TargetApiPackageName":    targetModuleName,
			"TargetApiPackageVersion": moduleVersion,

			"SourceApiPackageName":    sourceModuleName,
			"SourceApiPackageVersion": moduleVersion,

			"Service": serviceName,
			"Model":   name,
			"IsGrpc":  isGrpcService,
			"UseRepo": useRepo,
		},
	}

	if _, err := g.GenerateService(context.Background(), o); err != nil {
		return err
	}

	return nil
}

func WriteServerPackageCode(
	outputPath string,
	projectName string,
	serviceType string,
	serviceName string,
	services map[string]string,
) error {
	g := generators.NewGoGenerator()

	switch strings.ToLower(serviceType) {
	case "grpc":
		o := code_generator.Options{
			OutDir: outputPath,
			Module: projectName,
			Vars: map[string]interface{}{
				"Service":  serviceName,
				"Services": services,
			},
		}
		if _, err := g.GenerateGrpcServer(context.Background(), o); err != nil {
			return err
		}

	case "rest":
		o := code_generator.Options{
			OutDir: outputPath,
			Module: projectName,
			Vars: map[string]interface{}{
				"Service":  serviceName,
				"Services": services,
			},
		}
		if _, err := g.GenerateRestServer(context.Background(), o); err != nil {
			return err
		}

	default:
		return errors.New("sqlproto: unsupported service type: " + serviceType)
	}

	return nil
}

func WriteInitWireCode(
	outputPath string,

	packageName string,
	postfix string,
	services []string,
) error {
	g := generators.NewGoGenerator()

	var newFunctions []string
	for _, service := range services {
		newFunction := "New" + stringcase.UpperCamelCase(service) + postfix
		newFunctions = append(newFunctions, newFunction)
	}

	opts := code_generator.Options{
		OutDir: outputPath,
		Vars: map[string]interface{}{
			"Package":      packageName,
			"NewFunctions": newFunctions,
		},
	}
	_, err := g.GenerateInit(context.Background(), opts)
	return err
}

func WriteWireCode(
	outputPath string,

	projectName string,
	serviceName string,
) error {
	g := generators.NewGoGenerator()
	opts := code_generator.Options{
		OutDir: outputPath,
		Module: projectName,
		Vars: map[string]interface{}{
			"Service": serviceName,
		},
	}
	_, err := g.GenerateWire(context.Background(), opts)
	return err
}

func WriteMainCode(
	outputPath string,

	projectName string,
	serviceName string,

	servers []string,
) error {
	g := generators.NewGoGenerator()

	opts := code_generator.Options{
		OutDir: outputPath,
		Module: projectName,
		Vars: map[string]interface{}{
			"Service":                  serviceName,
			"ServerImports":            generators.ServerImportPaths(servers),
			"ServerFormalParameters":   generators.ServerFormalParameters(servers),
			"ServerTransferParameters": generators.ServerTransferParameters(servers),
		},
	}

	_, err := g.GenerateMain(context.Background(), opts)
	return err
}
