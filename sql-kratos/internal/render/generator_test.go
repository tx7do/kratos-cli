package render

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/jinzhu/inflection"
)

func TestPlural(t *testing.T) {
	t.Log(inflection.Plural("user"))
	t.Log(inflection.Plural("dict"))
	t.Log(inflection.Plural("admin_login_log"))
	t.Log(inflection.Plural("admin-login-log"))
	t.Log(inflection.Plural("adminLoginLog"))
}

func TestSingular(t *testing.T) {
	t.Log(inflection.Singular("users"))
	t.Log(inflection.Singular("dicts"))
	t.Log(inflection.Singular("admin_login_logs"))
	t.Log(inflection.Singular("admin-login-logs"))
	t.Log(inflection.Singular("adminLoginLogs"))
}

func TestWriteInitWireCode(t *testing.T) {
	var err error

	serviceInit := InitWireTemplateData{
		Package: "service",
		Postfix: "Service",
		ServiceNames: []string{
			"User",
			"Tenant",
		},
	}
	err = WriteInitWireCode("./app/user/internal/", serviceInit)
	assert.Nil(t, err)

	dataInit := InitWireTemplateData{
		Package: "data",
		Postfix: "Repo",
		ServiceNames: []string{
			"User",
			"Tenant",
		},
	}
	err = WriteInitWireCode("./app/user/internal/", dataInit)
	assert.Nil(t, err)

	serverInit := InitWireTemplateData{
		Package: "server",
		Postfix: "Server",
		ServiceNames: []string{
			"Service",
			"GRPC",
		},
	}
	err = WriteInitWireCode("./app/user/service/internal/", serverInit)
	assert.Nil(t, err)
}

func TestWriteWireCode(t *testing.T) {
	var err error

	serviceInit := WireTemplateData{
		Project: "kratos-admin",
		Service: "user",
	}
	err = WriteWireCode("./app/user/service/cmd/server", serviceInit)
	assert.Nil(t, err)
}

func TestWriteMainCode(t *testing.T) {
	var err error

	serviceInit := MainTemplateData{
		Project: "kratos-admin",
		Service: "user",
		Servers: []string{"REST", "GRPC"},
	}
	err = WriteMainCode("./app/user/service/cmd/server", serviceInit)
	assert.Nil(t, err)
}

func TestWriteServicePackageCode(t *testing.T) {
	var err error

	dataUser := ServiceTemplateData{
		Project: "kratos-admin",

		Service: "admin",
		Name:    "user",

		SourceModuleName: "user",
		TargetModuleName: "admin",
		Version:          "v1",

		UseRepo: false,
		IsGrpc:  false,
	}
	err = WriteServicePackageCode("./app/user/service/internal/", dataUser)
	assert.Nil(t, err)

	dataTenant := ServiceTemplateData{
		Project: "kratos-admin",

		Service: "user",
		Name:    "tenant",

		SourceModuleName: "user",
		TargetModuleName: "user",
		Version:          "v1",

		UseRepo: true,
		IsGrpc:  true,
	}
	err = WriteServicePackageCode("./app/user/service/internal/", dataTenant)
	assert.Nil(t, err)
}

func TestWriteEntDataPackageCode(t *testing.T) {
	var err error

	data := DataTemplateData{
		Project: "kratos-admin",

		Service: "user",
		Name:    "user",

		Module:  "user",
		Version: "v1",

		Fields: DataFieldArray{
			{Name: "id", Type: "int64", Comment: "用户ID"},
			{Name: "last_login_time", Type: "int64", Comment: "最后登录时间"},
		},
	}

	err = WriteEntDataPackageCode("./app/user/service/internal/", data)
	assert.Nil(t, err)
}

func TestWriteServerCode(t *testing.T) {
	var err error

	dataGrpc := ServerTemplateData{
		Project: "kratos-admin",
		Type:    "grpc",
		Service: "user",
		Services: map[string]string{
			"User":   "user",
			"Tenant": "user",
		},
	}
	err = WriteServerPackageCode("./app/user/service/internal/", dataGrpc)
	assert.Nil(t, err)

	dataRest := ServerTemplateData{
		Project: "kratos-admin",
		Type:    "rest",
		Service: "admin",
		Services: map[string]string{
			"User":   "user",
			"Tenant": "user",
		},
	}
	err = WriteServerPackageCode("./app/user/service/internal/", dataRest)
	assert.Nil(t, err)
}
