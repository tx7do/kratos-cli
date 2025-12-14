package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"

	authzEngine "github.com/tx7do/kratos-authz/engine"
	authz "github.com/tx7do/kratos-authz/middleware"

	swaggerUI "github.com/tx7do/kratos-swagger-ui"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"{{.Project}}/app/{{.Service}}/service/cmd/server/assets"
	"{{.Project}}/app/{{.Service}}/service/internal/service"

	{{lower .Service}}V1 "{{.Project}}/api/gen/go/{{.Service}}/service/v1"

	"{{.Project}}/pkg/middleware/auth"
	applogging "{{.Project}}/pkg/middleware/logging"
)

// NewWhiteListMatcher 创建jwt白名单
func newRestWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]bool)
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewMiddleware 创建中间件
func newRestMiddleware(
	logger log.Logger,
	authenticator authnEngine.Authenticator,
	authorizer authzEngine.Engine,
) []middleware.Middleware {
	var ms []middleware.Middleware

	ms = append(ms, logging.Server(logger))

	ms = append(ms, selector.Server(
		authn.Server(authenticator),
		auth.Server(),
		authz.Server(authorizer),
	).Match(newRestWhiteListMatcher()).Build())

	return ms
}

// NewRestServer new an HTTP server.
func NewRestServer(
	cfg *conf.Bootstrap, logger log.Logger,

	authenticator authnEngine.Authenticator,
	authorizer authzEngine.Engine,
{{range $key, $value := .Services}}
    {{lower $key}}Service *service.{{pascal $key}}Service,
{{- end}}
) *http.Server {
	if cfg == nil || cfg.Server == nil || cfg.Server.Rest == nil {
		return nil
	}

	srv := rpc.CreateRestServer(cfg, newRestMiddleware(logger, authenticator, authorizer)...)
{{range $key, $value := .Services}}
    {{lower $value}}V1.Register{{pascal $key}}ServiceHTTPServer(srv, {{lower $key}}Service)
{{- end}}

	if cfg.GetServer().GetRest().GetEnableSwagger() {
		swaggerUI.RegisterSwaggerUIServerWithOption(
			srv,
			swaggerUI.WithTitle("{{pascal .Service}} Service"),
			swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
		)
	}

	return srv
}
