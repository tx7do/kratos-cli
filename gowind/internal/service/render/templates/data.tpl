package data

import (
	"github.com/go-kratos/kratos/v2/log"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
)

// Data .
type Data struct {
	log *log.Helper
}

// NewData .
func NewData(
	logger log.Logger,
) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "data/{{.Service}}-service"))

	d := &Data{
		log: l,
	}

	return d, func() {
		l.Info("closing the data resources")
	}, nil
}