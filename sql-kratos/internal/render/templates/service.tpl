package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
{{if not .IsGrpc}}
	"github.com/tx7do/go-utils/trans"
{{end}}
	"google.golang.org/protobuf/types/known/emptypb"
	pagination "github.com/tx7do/kratos-bootstrap/api/gen/go/pagination/v1"

	"{{.Project}}/app/{{.Service}}/service/internal/data"

	{{.TargetApiPackage}} "{{.Project}}/api/gen/go/{{lower .TargetModuleName}}/service/v1"
{{if not .IsSameApi}}	{{.SourceApiPackage}} "{{.Project}}/api/gen/go/{{lower .SourceModuleName}}/service/v1"{{end}}

{{if not .IsGrpc}}
	"{{.Project}}/pkg/middleware/auth"
{{end -}}
)

type {{.ClassName}} struct {
	{{.ServiceInterface}}

	log *log.Helper

	{{.DataSourceVar}} {{.DataSourceType}}
}

func New{{.ClassName}}(
    logger log.Logger,
    {{.DataSourceVar}} {{.DataSourceType}},
) *{{.ClassName}} {
	return &{{.ClassName}}{
		log: log.NewHelper(log.With(logger, "module", "{{.LowerName}}/service/{{.Service}}-service")),
		{{.DataSourceVar}}:  {{.DataSourceVar}},
	}
}

func (s *{{.ClassName}}) List(ctx context.Context, req *pagination.PagingRequest) (*{{.SourceApiPackage}}.List{{.PascalName}}Response, error) {
	return s.{{.DataSourceVar}}.List(ctx, req)
}

func (s *{{.ClassName}}) Get(ctx context.Context, req *{{.SourceApiPackage}}.Get{{.PascalName}}Request) (*{{.SourceApiPackage}}.{{.PascalName}}, error) {
	return s.{{.DataSourceVar}}.Get(ctx, req)
}

func (s *{{.ClassName}}) Create(ctx context.Context, req *{{.SourceApiPackage}}.Create{{.PascalName}}Request) (*{{.SourceApiPackage}}.{{.PascalName}}, error) {
	if req == nil || req.Data == nil {
		return nil, {{.SourceApiPackage}}.ErrorBadRequest("invalid parameter")
	}

{{if not .IsGrpc}}
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	req.Data.CreateBy = trans.Ptr(operator.UserId)
{{end}}

{{if not .UseRepo}}
	if result, err := s.{{.DataSourceVar}}.Create(ctx, req); err != nil {
		return nil, err
	} else {
	    return result, nil
	}
{{else}}
    return s.{{.DataSourceVar}}.Create(ctx, req)
{{end -}}
}

func (s *{{.ClassName}}) Update(ctx context.Context, req *{{.SourceApiPackage}}.Update{{.PascalName}}Request) (*emptypb.Empty, error) {
	if req == nil || req.Data == nil {
		return nil, {{.SourceApiPackage}}.ErrorBadRequest("invalid parameter")
	}

{{if not .IsGrpc}}
	// 获取操作人信息
	operator, err := auth.FromContext(ctx)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	req.Data.UpdateBy = trans.Ptr(operator.UserId)
{{end -}}

{{if not .UseRepo}}
	if _, err := s.{{.DataSourceVar}}.Update(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
{{else}}
	if _, err := s.{{.DataSourceVar}}.Update(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
{{end -}}
}

func (s *{{.ClassName}}) Delete(ctx context.Context, req *{{.SourceApiPackage}}.Delete{{.PascalName}}Request) (*emptypb.Empty, error) {
    if req == nil {
        return nil, {{.SourceApiPackage}}.ErrorBadRequest("invalid parameter")
    }

{{if not .UseRepo}}
	if err := s.{{.DataSourceVar}}.Delete(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
{{else}}
    return s.{{.DataSourceVar}}.Delete(ctx, req)
{{end -}}
}
