package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"{{.Module}}/app/{{lower $.Service}}/service/internal/service"
)

// NewGrpcServer new a gRPC server.
func NewGrpcServer(
	cfg *conf.Bootstrap, logger log.Logger,
) *grpc.Server {
	srv := rpc.CreateGrpcServer(cfg, logging.Server(logger))

	return srv
}
