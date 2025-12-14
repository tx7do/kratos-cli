package data

import (
	"context"
	"errors"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/tx7do/go-utils/copierutil"
	"github.com/tx7do/go-utils/crypto"
	"github.com/tx7do/go-utils/mapper"

	pagination "github.com/tx7do/go-crud/api/gen/go/pagination/v1"
	entCurd "github.com/tx7do/go-crud/entgo"

	"{{.Project}}/app/{{.Service}}/service/internal/data/ent"
	"{{.Project}}/app/{{.Service}}/service/internal/data/ent/predicate"
	"{{.Project}}/app/{{.Service}}/service/internal/data/ent/{{.LowerName}}"

	{{.ApiPackage}} "{{.Project}}/api/gen/go/{{.Module}}/service/v1"
)

type {{.ClassName}} struct {
	data         *Data
	log          *log.Helper

	mapper     *mapper.CopierMapper[{{.ApiPackage}}.{{.PascalName}}, ent.{{.PascalName}}]
	repository *entCurd.Repository[
		ent.{{.PascalName}}Query, ent.{{.PascalName}}Select,
		ent.{{.PascalName}}Create, ent.{{.PascalName}}CreateBulk,
		ent.{{.PascalName}}Update, ent.{{.PascalName}}UpdateOne,
		ent.{{.PascalName}}Delete,
		predicate.{{.PascalName}},
		{{.ApiPackage}}.{{.PascalName}}, ent.{{.PascalName}},
	]
}

func New{{.ClassName}}(data *Data, logger log.Logger) *{{.ClassName}} {
	repo := &{{.ClassName}}{
		log:  log.NewHelper(log.With(logger, "module", "{{.LowerName}}/repo/{{.Service}}-service")),
		data: data,
		mapper: mapper.NewCopierMapper[{{.ApiPackage}}.{{.PascalName}}, ent.{{.PascalName}}](),
	}

	repo.init()

	return repo
}

func (r *{{.ClassName}}) init() {
	r.repository = entCurd.NewRepository[
		ent.{{.PascalName}}Query, ent.{{.PascalName}}Select,
		ent.{{.PascalName}}Create, ent.{{.PascalName}}CreateBulk,
		ent.{{.PascalName}}Update, ent.{{.PascalName}}UpdateOne,
		ent.{{.PascalName}}Delete,
		predicate.{{.PascalName}},
		userV1.{{.PascalName}}, ent.{{.PascalName}},
	](r.mapper)

	r.mapper.AppendConverters(copierutil.NewTimeStringConverterPair())
	r.mapper.AppendConverters(copierutil.NewTimeTimestamppbConverterPair())
}

func (r *{{.ClassName}}) Count(ctx context.Context, whereCond []func(s *sql.Selector)) (int, error) {
	builder := r.data.db.Client().{{.PascalName}}.Query()
	if len(whereCond) != 0 {
		builder.Modify(whereCond...)
	}

	count, err := builder.Count(ctx)
	if err != nil {
		r.log.Errorf("query count failed: %s", err.Error())
		return 0, {{.ApiPackage}}.ErrorInternalServerError("query count failed")
	}

	return count, nil
}

func (r *{{.ClassName}}) List(ctx context.Context, req *pagination.PagingRequest) (*{{.ApiPackage}}.List{{.PascalName}}Response, error) {
	if req == nil {
		return nil, {{.ApiPackage}}.ErrorBadRequest("invalid parameter")
	}

	builder := r.data.db.Client().{{.PascalName}}.Query()

    ret, err := r.repository.ListWithPaging(ctx, builder, builder.Clone(), req)
    if err != nil {
        return nil, err
    }
    if ret == nil {
        return &{{.ApiPackage}}.List{{.PascalName}}Response{Total: 0, Items: nil}, nil
    }

    return &{{.ApiPackage}}.List{{.PascalName}}Response{
        Total: ret.Total,
        Items: ret.Items,
    }, nil
}

func (r *{{.ClassName}}) IsExist(ctx context.Context, id uint32) (bool, error) {
	exist, err := r.data.db.Client().{{.PascalName}}.Query().
		Where({{.LowerName}}.IDEQ(id)).
		Exist(ctx)
	if err != nil {
		r.log.Errorf("query exist failed: %s", err.Error())
		return false, {{.ApiPackage}}.ErrorInternalServerError("query exist failed")
	}
	return exist, nil
}

func (r *{{.ClassName}}) Get(ctx context.Context, req *{{.ApiPackage}}.Get{{.PascalName}}Request) (*{{.ApiPackage}}.{{.PascalName}}, error) {
	if req == nil {
		return nil, {{.ApiPackage}}.ErrorBadRequest("invalid parameter")
	}

    var whereCond []func(s *sql.Selector)
    switch req.QueryBy.(type) {
    case *{{.ApiPackage}}.GetUserRequest_Id:
        whereCond = append(whereCond, {{.LowerName}}.IDEQ(req.GetId()))
    default:
        whereCond = append(whereCond, {{.LowerName}}.IDEQ(req.GetId()))
    }

    builder := r.data.db.Client().{{.PascalName}}.Query()
	dto, err := r.repository.Get(ctx, builder, req.GetViewMask(), whereCond...)
	if err != nil {
		return nil, err
	}

	return dto, err
}

func (r *{{.ClassName}}) Create(ctx context.Context, req *{{.ApiPackage}}.Create{{.PascalName}}Request) error {
	if req == nil || req.Data == nil {
		return {{.ApiPackage}}.ErrorBadRequest("invalid parameter")
	}

	builder := r.data.db.Client().{{.PascalName}}.Create()

	builder{{range .Fields}}.{{newline}}		{{.EntSetNillableFunc}}
{{- end}}

    builder.SetNillableCreateBy(req.Data.CreateBy)

	if req.Data.CreateTime == nil {
		builder.SetCreateTime(time.Now())
	} else {
		builder.SetNillableCreateTime(timeutil.StringTimeToTime(req.Data.CreateTime))
	}

	if req.Data.Id != nil {
		builder.SetID(req.Data.GetId())
	}

	if err := builder.Exec(ctx); err != nil {
		r.log.Errorf("insert one data failed: %s", err.Error())
		return {{.ApiPackage}}.ErrorInternalServerError("insert data failed")
	}

	return nil
}

func (r *{{.ClassName}}) Update(ctx context.Context, req *{{.ApiPackage}}.Update{{.PascalName}}Request) error {
	if req == nil || req.Data == nil {
		return {{.ApiPackage}}.ErrorBadRequest("invalid parameter")
	}

	// 如果不存在则创建
	if req.GetAllowMissing() {
		exist, err := r.IsExist(ctx, req.GetData().GetId())
		if err != nil {
			return err
		}
		if !exist {
			createReq := &{{.ApiPackage}}.Create{{.PascalName}}Request{Data: req.Data}
			createReq.Data.CreateBy = createReq.Data.UpdateBy
			createReq.Data.UpdateBy = nil
			return r.Create(ctx, createReq)
		}
	}

	builder := r.data.db.Client().{{.PascalName}}.UpdateOneID(req.Data.GetId())
	result, err := r.repository.UpdateOne(ctx, builder, req.Data, req.GetUpdateMask(),
		func(dto *{{.ApiPackage}}.{{.PascalName}}) {
			builder.
				SetUpdatedAt(time.Now())
		},
	)

	return result, err
}

func (r *{{.ClassName}}) Delete(ctx context.Context, req *{{.ApiPackage}}.Delete{{.PascalName}}Request) error {
	if req == nil {
		return {{.ApiPackage}}.ErrorBadRequest("invalid parameter")
	}

	builder := r.data.db.Client().{{.PascalName}}.Delete()
	_, err := r.repository.Delete(ctx, builder, func(s *sql.Selector) {
		s.Where(sql.EQ({{.LowerName}}.FieldID, req.GetId()))
	})

	return nil
}
