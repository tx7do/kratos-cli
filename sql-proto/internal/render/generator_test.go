package render

import (
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

func TestWriteGrpcServiceProto(t *testing.T) {
	data := GrpcProtoTemplateData{
		Module:  "user",
		Name:    "user",
		Comment: "用户",
		Version: "v1",

		Fields: []ProtoField{
			{Name: "id", Type: "int64", Number: 1, Comment: "用户ID"},
			{Name: "last_login_time", Type: "int64", Number: 2, Comment: "最后登录时间"},
		},
	}

	WriteGrpcServiceProto("./api/protos", data)
}

func TestWriteRestServiceProto(t *testing.T) {
	data := RestProtoTemplateData{
		SourceModule: "user",
		TargetModule: "admin",
		Name:         "user",
		Version:      "v1",
		Comment:      "用户",
	}

	WriteRestServiceProto("./api/protos", data)
}

func TestWriteInitWireCode(t *testing.T) {
	serviceInit := InitWireTemplateData{
		Package: "service",
		Postfix: "Service",
		ServiceNames: []string{
			"User",
			"Tenant",
		},
	}
	WriteInitWireCode("./app/user/internal/", serviceInit)

	dataInit := InitWireTemplateData{
		Package: "data",
		Postfix: "Repo",
		ServiceNames: []string{
			"User",
			"Tenant",
		},
	}
	WriteInitWireCode("./app/user/internal/", dataInit)

	serverInit := InitWireTemplateData{
		Package: "server",
		Postfix: "Server",
		ServiceNames: []string{
			"Service",
			"GRPC",
		},
	}
	WriteInitWireCode("./app/user/service/internal/", serverInit)
}

func TestWriteWireCode(t *testing.T) {
	serviceInit := WireTemplateData{
		Project: "kratos-admin",
		Service: "user",
	}
	WriteWireCode("./app/user/service/cmd/server", serviceInit)
}

func TestWriteMainCode(t *testing.T) {
	serviceInit := MainTemplateData{
		Project:     "kratos-admin",
		Service:     "user",
		EnableREST:  true,
		EnableGRPC:  true,
		EnableAsynq: false,
		EnableSSE:   false,
		EnableKafka: false,
	}
	WriteMainCode("./app/user/service/cmd/server", serviceInit)
}

func TestWriteRestServicePackageCode(t *testing.T) {
	dataUser := ServiceTemplateData{
		Project: "kratos-admin",

		Service: "admin",
		Name:    "user",

		SourceApi: "user",
		TargetApi: "admin",
		Version:   "v1",

		UseRepo: false,
		IsGrpc:  false,
	}
	WriteGrpcServicePackageCode("./app/user/service/internal/", dataUser)

	dataTenant := ServiceTemplateData{
		Project: "kratos-admin",

		Service: "user",
		Name:    "tenant",

		SourceApi: "user",
		TargetApi: "user",
		Version:   "v1",

		UseRepo: true,
		IsGrpc:  true,
	}
	WriteGrpcServicePackageCode("./app/user/service/internal/", dataTenant)
}

func TestWriteEntDataPackageCode(t *testing.T) {
	data := DataTemplateData{
		Project: "kratos-admin",

		Service: "user",
		Name:    "user",

		Module:  "user",
		Version: "v1",

		Fields: []ProtoField{
			{Name: "id", Type: "int64", Number: 1, Comment: "用户ID"},
			{Name: "last_login_time", Type: "int64", Number: 2, Comment: "最后登录时间"},
		},
	}

	WriteEntDataPackageCode("./app/user/service/internal/", data)
}

func TestWriteServerCode(t *testing.T) {
	dataGrpc := ServerTemplateData{
		Project: "kratos-admin",
		Type:    "grpc",
		Service: "user",
		Services: map[string]string{
			"User":   "user",
			"Tenant": "user",
		},
	}
	WriteServerPackageCode("./app/user/service/internal/", dataGrpc)

	dataRest := ServerTemplateData{
		Project: "kratos-admin",
		Type:    "rest",
		Service: "admin",
		Services: map[string]string{
			"User":   "user",
			"Tenant": "user",
		},
	}
	WriteServerPackageCode("./app/user/service/internal/", dataRest)
}
