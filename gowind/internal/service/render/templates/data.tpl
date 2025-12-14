package data

import (
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
{{if .HasRedis}}
    "github.com/redis/go-redis/v9"
{{end}}{{if .HasGorm}}
    gormCrud "github.com/tx7do/go-crud/gorm"
{{end}}{{if .HasEnt}}
    entCrud "github.com/tx7do/go-crud/entgo"
    "{{.Module}}/app/{{.Service}}/service/internal/data/ent"
{{end}})

// Data .
type Data struct {
	log *log.Helper
{{if .HasRedis}}
    rdb *redis.Client
{{end}}{{if .HasGorm}}
    gorm *gormCrud.Client
{{end}}{{if .HasEnt}}
    db *entCrud.EntClient[*ent.Client]
{{end}}}

// NewData .
func NewData(
	logger log.Logger,
{{if .HasRedis}}
    rdb *redis.Client,
{{end}}{{if .HasGorm}}
    gorm *gormCrud.Client,
{{end}}{{if .HasEnt}}
    db *entCrud.EntClient[*ent.Client],
{{end}}) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "data/{{.Service}}-service"))

	d := &Data{
		log: l,
{{if .HasRedis}}
		rdb: rdb,
{{end}}{{if .HasGorm}}
		gorm: gorm,
{{end}}{{if .HasEnt}}
		db: db,
{{end}}	}

	return d, func() {
		l.Info("closing the data resources")
{{if .HasRedis}}
        if d.rdb != nil {
            if err := d.rdb.Close(); err != nil {
                l.Error(err)
            }
        }
{{end}}{{if .HasEnt}}
        if d.db != nil {
            if err := d.db.Close(); err != nil {
                l.Error(err)
            }
        }
{{end}}	}, nil
}
