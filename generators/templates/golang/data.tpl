package data

import (
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	"github.com/tx7do/kratos-bootstrap/bootstrap"
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
{{- if .HasRedis}}
    rdb *redis.Client
{{- end}}
{{- if .HasGorm}}
    gorm *gormCrud.Client
{{- end}}
{{- if .HasEnt}}
    db *entCrud.EntClient[*ent.Client]
{{- end}}
}

// NewData .
func NewData(
	ctx *bootstrap.Context,
{{- if .HasEnt}}
    db *entCrud.EntClient[*ent.Client],
{{- end}}
{{- if .HasGorm}}
    gorm *gormCrud.Client,
{{- end}}
{{- if .HasRedis}}
    rdb *redis.Client,
{{- end}}
) (*Data, func(), error) {
	d := &Data{
		log: ctx.NewLoggerHelper("data/{{.Service}}-service"),
{{- if .HasEnt }}
		db:   db,
{{- end }}
{{- if .HasGorm }}
		gorm: gorm,
{{- end }}
{{- if .HasRedis }}
		rdb:  rdb,
{{- end }}
	}

	cleanup := func() {
		d.log.Info("closing the data resources")
{{- if .HasEnt }}
		if d.db != nil {
			if err := d.db.Close(); err != nil {
				d.log.Error(err)
			}
		}
{{- end }}
{{- if .HasGorm }}
		if d.gorm != nil {
			if err := d.gorm.Close(); err != nil {
				d.log.Error(err)
			}
		}
{{- end }}
{{- if .HasRedis }}
		if d.rdb != nil {
			if err := d.rdb.Close(); err != nil {
				d.log.Error(err)
			}
		}
{{- end }}
	}

	return d, cleanup, nil
}
