package data

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"

	"github.com/tx7do/go-utils/copierutil"
	entgo "github.com/tx7do/go-utils/entgo/query"
	entgoUpdate "github.com/tx7do/go-utils/entgo/update"
	"github.com/tx7do/go-utils/fieldmaskutil"
	"github.com/tx7do/go-utils/timeutil"
	"github.com/tx7do/go-utils/trans"
	pagination "github.com/tx7do/kratos-bootstrap/api/gen/go/pagination/v1"

	"{{.Project}}/app/{{.Service}}/service/internal/data/ent"
	"{{.Project}}/app/{{.Service}}/service/internal/data/ent/{{.LowerName}}"

	{{.ApiPackage}} "{{.Project}}/api/gen/go/{{.Module}}/service/v1"
)

type {{.ClassName}} struct {
	data         *Data
	log          *log.Helper
	copierOption copier.Option
}

func New{{.ClassName}}(data *Data, logger log.Logger) *{{.ClassName}} {
	repo := &{{.ClassName}}{
		log:  log.NewHelper(log.With(logger, "module", "{{.LowerName}}/repo/{{.Service}}-service")),
		data: data,
	}

	repo.init()

	return repo
}

func (r *{{.ClassName}}) init() {
	r.copierOption = copier.Option{
		Converters: []copier.TypeConverter{},
	}

	r.copierOption.Converters = append(r.copierOption.Converters, copierutil.NewTimeStringConverterPair()...)
	r.copierOption.Converters = append(r.copierOption.Converters, copierutil.NewTimeTimestamppbConverterPair()...)
}

func (r *{{.ClassName}}) toProto(in *ent.{{.PascalName}}) *{{.ApiPackage}}.{{.PascalName}} {
	if in == nil {
		return nil
	}

	var out {{.ApiPackage}}.{{.PascalName}}
	_ = copier.CopyWithOption(&out, in, r.copierOption)

	return &out
}

func (r *{{.ClassName}}) toEnt(in *{{.ApiPackage}}.{{.PascalName}}) *ent.{{.PascalName}} {
	if in == nil {
		return nil
	}

	var out ent.{{.PascalName}}
	_ = copier.CopyWithOption(&out, in, r.copierOption)

	return &out
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

	err, whereSelectors, querySelectors := entgo.BuildQuerySelector(
		req.GetQuery(), req.GetOrQuery(),
		req.GetPage(), req.GetPageSize(), req.GetNoPaging(),
		req.GetOrderBy(), {{.LowerName}}.FieldCreateTime,
		req.GetFieldMask().GetPaths(),
	)
	if err != nil {
		r.log.Errorf("parse list param error [%s]", err.Error())
		return nil, {{.ApiPackage}}.ErrorBadRequest("invalid query parameter")
	}

	if querySelectors != nil {
		builder.Modify(querySelectors...)
	}

	results, err := builder.All(ctx)
	if err != nil {
		r.log.Errorf("query list failed: %s", err.Error())
		return nil, {{.ApiPackage}}.ErrorInternalServerError("query list failed")
	}

	items := make([]*{{.ApiPackage}}.{{.PascalName}}, 0, len(results))
	for _, res := range results {
		item := r.toProto(res)
		items = append(items, item)
	}

	count, err := r.Count(ctx, whereSelectors)
	if err != nil {
		return nil, err
	}

	return &{{.ApiPackage}}.List{{.PascalName}}Response{
		Total: uint32(count),
		Items: items,
	}, err
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

	ret, err := r.data.db.Client().{{.PascalName}}.Get(ctx, req.GetId())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, {{.ApiPackage}}.Error{{.PascalName}}NotFound("{{.LowerName}} not found")
		}

		r.log.Errorf("query one data failed: %s", err.Error())

		return nil, {{.ApiPackage}}.ErrorInternalServerError("query data failed")
	}

	return r.toProto(ret), nil
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

	if req.UpdateMask != nil {
		req.UpdateMask.Normalize()
		if !req.UpdateMask.IsValid(req.Data) {
			r.log.Errorf("invalid field mask [%v]", req.UpdateMask)
			return {{.ApiPackage}}.ErrorBadRequest("invalid field mask")
		}
		fieldmaskutil.Filter(req.GetData(), req.UpdateMask.GetPaths())
	}

	builder := r.data.db.Client().{{.PascalName}}.UpdateOneID(req.Data.GetId())

	builder{{range .Fields}}.{{newline}}		{{.EntSetNillableFunc}}
{{- end}}

    builder.SetNillableUpdateBy(req.Data.UpdateBy)

	if req.Data.UpdateTime == nil {
		builder.SetUpdateTime(time.Now())
	} else {
		builder.SetNillableUpdateTime(timeutil.StringTimeToTime(req.Data.UpdateTime))
	}

	if req.UpdateMask != nil {
		nilPaths := fieldmaskutil.NilValuePaths(req.Data, req.GetUpdateMask().GetPaths())
		nilUpdater := entgoUpdate.BuildSetNullUpdater(nilPaths)
		if nilUpdater != nil {
			builder.Modify(nilUpdater)
		}
	}

	if err := builder.Exec(ctx); err != nil {
		r.log.Errorf("update one data failed: %s", err.Error())
		return {{.ApiPackage}}.ErrorInternalServerError("update data failed")
	}

	return nil
}

func (r *{{.ClassName}}) Delete(ctx context.Context, req *{{.ApiPackage}}.Delete{{.PascalName}}Request) error {
	if req == nil {
		return {{.ApiPackage}}.ErrorBadRequest("invalid parameter")
	}

	if err := r.data.db.Client().{{.PascalName}}.DeleteOneID(req.GetId()).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			return {{.ApiPackage}}.ErrorNotFound("{{.LowerName}} not found")
		}

		r.log.Errorf("delete one data failed: %s", err.Error())

		return {{.ApiPackage}}.ErrorInternalServerError("delete failed")
	}

	return nil
}
