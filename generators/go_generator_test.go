package generators

import (
	"context"
	"testing"

	"github.com/tx7do/go-utils/code_generator"
)

func TestGoGenerator_Template_Main(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir:      "./output",
		ProjectName: "MyProject",
		Vars: map[string]interface{}{
			"Service":                  "user",
			"ServerImports":            []string{"github.com/example/myproject/server"},
			"ServerFormalParameters":   []string{"hs http.Server"},
			"ServerTransferParameters": []string{"hs"},
		},
	}

	if _, err := g.GenerateMain(context.Background(), opts); err != nil {
		t.Fatalf("Generate main.go failed: %v", err)
	}
}

func TestGoGenerator_Template_Wire(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir:      "./output",
		ProjectName: "MyProject",
		Module:      "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Service": "user",
		},
	}

	if _, err := g.GenerateWire(context.Background(), opts); err != nil {
		t.Fatalf("Generate wire.go failed: %v", err)
	}
}

func TestGoGenerator_Template_Init(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir:      "./output",
		ProjectName: "MyProject",
		Module:      "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Package":      "data",
			"NewFunctions": []string{"NewUserRepo", "NewOrderRepo"},
		},
	}

	if _, err := g.GenerateInit(context.Background(), opts); err != nil {
		t.Fatalf("Generate init.go failed: %v", err)
	}
}

func TestGoGenerator_Template_Data(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Service":  "user",
			"HasRedis": true,
			"HasGorm":  true,
			"HasEnt":   true,
		},
	}

	if _, err := g.GenerateData(context.Background(), opts); err != nil {
		t.Fatalf("Generate data.go failed: %v", err)
	}
}

func TestGoGenerator_Template_EntClient(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Service": "user",
		},
	}

	if _, err := g.GenerateEntClient(context.Background(), opts); err != nil {
		t.Fatalf("Generate ent_client.go failed: %v", err)
	}
}

func TestGoGenerator_Template_EntRepo(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Service":    "user",
			"ApiPackage": "userV1",
			"Model":      "user",
		},
	}

	if _, err := g.GenerateEntRepo(context.Background(), opts); err != nil {
		t.Fatalf("Generate ent_repo.go failed: %v", err)
	}
}

func TestGoGenerator_Template_GormClient(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Service": "user",
		},
	}

	if _, err := g.GenerateGormClient(context.Background(), opts); err != nil {
		t.Fatalf("Generate gorm_client.go failed: %v", err)
	}
}

func TestGoGenerator_Template_GormInit(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir:     "./output",
		OutputName: "gorm_init.go",
		Module:     "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Service": "user",
		},
	}

	if _, err := g.GenerateGormInit(context.Background(), opts); err != nil {
		t.Fatalf("Generate gorm_init.go failed: %v", err)
	}
}

func TestGoGenerator_Template_GormRepo(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Service":    "user",
			"ApiPackage": "userV1",
			"Model":      "user",
		},
	}

	if _, err := g.GenerateGormRepo(context.Background(), opts); err != nil {
		t.Fatalf("Generate gorm_repo.go failed: %v", err)
	}
}

func TestGoGenerator_Template_GrpcServiceProto(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Package":   "user.service.v1",
			"Model":     "user",
			"ModelName": "用户",
			"Fields": []ProtoFieldData{
				{Name: "id", Type: "int64", Comment: "用户ID", Number: 1},
				{Name: "name", Type: "string", Comment: "用户名", Number: 2},
				{Name: "email", Type: "string", Comment: "用户邮箱", Number: 3},
			},
		},
	}

	if _, err := g.GenerateGrpcServiceProto(context.Background(), opts); err != nil {
		t.Fatalf("Generate grpc_proto.go failed: %v", err)
	}
}

func TestGoGenerator_Template_GrpcServer(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Service":  "user",
			"Packages": []string{"user"},
			"Services": map[string]string{"user": "userV1", "role": "userV1"},
		},
	}

	if _, err := g.GenerateGrpcServer(context.Background(), opts); err != nil {
		t.Fatalf("Generate grpc_server.go failed: %v", err)
	}
}

func TestGoGenerator_Template_RedisClient(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Service": "user",
		},
	}

	if _, err := g.GenerateRedisClient(context.Background(), opts); err != nil {
		t.Fatalf("Generate redis_client.go failed: %v", err)
	}
}

func TestGoGenerator_Template_RestServiceProto(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"TargetPackage": "admin.service.v1",
			"SourcePackage": "user.service.v1",
			"SourceProto":   "user/service/v1/user.proto",
			"Model":         "user",
			"Path":          "/admin/v1/users",
			"ModelName":     "用户",
		},
	}

	if _, err := g.GenerateRestServiceProto(context.Background(), opts); err != nil {
		t.Fatalf("Generate rest_proto.go failed: %v", err)
	}
}

func TestGoGenerator_Template_RestServer(t *testing.T) {
	g := NewGoGenerator()

	opts := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"Service":  "admin",
			"Services": map[string]string{"user": "userV1", "role": "userV1"},
		},
	}

	if _, err := g.GenerateRestServer(context.Background(), opts); err != nil {
		t.Fatalf("Generate rest_server.go failed: %v", err)
	}
}

func TestGoGenerator_Template_Service(t *testing.T) {
	g := NewGoGenerator()

	opts1 := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"TargetApiPackageName":    "user",
			"TargetApiPackageVersion": "v1",

			"SourceApiPackageName":    "user",
			"SourceApiPackageVersion": "v1",

			"Service": "user",
			"Model":   "user",
			"IsGrpc":  true,
		},
	}

	if _, err := g.GenerateService(context.Background(), opts1); err != nil {
		t.Fatalf("Generate service.go failed: %v", err)
	}

	opts2 := code_generator.Options{
		OutDir: "./output",
		Module: "github.com/example/myproject",
		Vars: map[string]interface{}{
			"TargetApiPackageName":    "admin",
			"TargetApiPackageVersion": "v1",

			"SourceApiPackageName":    "user",
			"SourceApiPackageVersion": "v1",

			"Service": "admin",
			"Model":   "role",
			"IsGrpc":  false,
		},
	}

	if _, err := g.GenerateService(context.Background(), opts2); err != nil {
		t.Fatalf("Generate service.go failed: %v", err)
	}
}
