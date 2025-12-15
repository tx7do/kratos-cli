package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"{{.Module}}/app/{{lower .Service}}/service/internal/service"
{{range $key, $value := .Packages}}
    {{lower $value}} "{{lower $.Module}}/api/gen/go/{{lower $value}}/service/v1"
{{- end}}
)

// NewGrpcServer creates a gRPC server.
func NewGrpcServer(
	cfg *conf.Bootstrap, logger log.Logger,
{{range $key, $value := .Services}}
    {{lower $key}}Service *service.{{pascal $key}}Service,
{{- end}}
) *grpc.Server {
	if cfg == nil || cfg.Server == nil || cfg.Server.Grpc == nil {
		return nil
	}

	srv := rpc.CreateGrpcServer(
		cfg,
		logging.Server(logger),
	)
{{range $key, $value := .Services}}
    {{lower $value}}V1.Register{{pascal $key}}ServiceServer(srv, {{lower $key}}Service)
{{- end}}

	return srv
}
