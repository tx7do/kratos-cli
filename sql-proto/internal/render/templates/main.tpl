package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"

{{if .EnableREST}}	"github.com/go-kratos/kratos/v2/transport/http"{{end}}
{{if .EnableGRPC}}	"github.com/go-kratos/kratos/v2/transport/grpc"{{end}}
{{if .EnableAsynq}}	"github.com/tx7do/kratos-transport/transport/asynq"{{end}}
{{if .EnableSSE}}	"github.com/tx7do/kratos-transport/transport/sse"{{end}}
{{if .EnableKafka}}	"github.com/tx7do/kratos-transport/transport/kafka"{{end}}
{{if .EnableMQTT}}	"github.com/tx7do/kratos-transport/transport/mqtt"{{end}}

	"github.com/tx7do/go-utils/trans"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"kratos-admin/pkg/service"
)

var version string

// go build -ldflags "-X main.version=x.y.z"

func newApp(
	lg log.Logger,
	re registry.Registrar,
{{if .EnableREST}}	hs *http.Server,{{end}}
{{if .EnableGRPC}}	gs *grpc.Server,{{end}}
{{if .EnableAsynq}}	as *asynq.Server,{{end}}
{{if .EnableSSE}}	ss *sse.Server,{{end}}
{{if .EnableKafka}}	ks *kafka.Server,{{end}}
{{if .EnableMQTT}}	ms *mqtt.Server,{{end}}
) *kratos.App {
	return bootstrap.NewApp(
		lg,
		re,
{{if .EnableREST}}		hs,{{end}}
{{if .EnableGRPC}}		gs,{{end}}
{{if .EnableAsynq}}		as,{{end}}
{{if .EnableSSE}}		ss,{{end}}
{{if .EnableKafka}}		ks,{{end}}
{{if .EnableMQTT}}		ms,{{end}}
	)
}

func main() {
	bootstrap.Bootstrap(initApp, trans.Ptr(service.{{.Service}}Service), trans.Ptr(version))
}
