package render

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
