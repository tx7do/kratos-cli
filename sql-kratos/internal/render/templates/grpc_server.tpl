package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"{{.Project}}/app/{{.Service}}/service/internal/service"
{{range $key, $value := .Modules}}
    {{lower $value}}V1 "{{lower $.Project}}/api/gen/go/{{lower $value}}/service/v1"
{{- end}}
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
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
