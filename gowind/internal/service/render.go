package service

import (
	"github.com/tx7do/kratos-cli/gowind/internal/service/render"
)

func WriteServerPackageCode(
	outputPath string,
	projectModule string,
	projectName string,
	serviceType string,
	serviceName string,
) error {
	data := render.ServerTemplateData{
		Module:  projectModule,
		Project: projectName,
		Type:    serviceType,
		Service: serviceName,
	}
	return render.WriteServerPackageCode(outputPath, data)
}

func WriteDataPackageCode(
	outputPath string,
	projectModule string,
	projectName string,
	serviceName string,
	dbClients []string,
) error {
	data := render.DataTemplateData{
		Module:    projectModule,
		Project:   projectName,
		Service:   serviceName,
		DBClients: dbClients,
	}
	return render.WriteDataPackageCode(outputPath, data)
}

func WriteInitWireCode(
	outputPath string,
	packageName string,
	postfix string,
	services []string,
) error {
	data := render.InitWireTemplateData{
		Package:      packageName,
		Postfix:      postfix,
		ServiceNames: services,
	}
	return render.WriteInitWireCode(outputPath, data)
}

func WriteWireCode(
	outputPath string,
	projectName string,
	serviceName string,
) error {
	data := render.WireTemplateData{
		Project: projectName,
		Service: serviceName,
	}
	return render.WriteWireCode(outputPath, data)
}

func WriteMainCode(
	outputPath string,
	projectName string,
	serviceName string,
	servers []string,
) error {
	data := render.MainTemplateData{
		Project: projectName,
		Service: serviceName,
		Servers: servers,
	}
	return render.WriteMainCode(outputPath, data)
}
