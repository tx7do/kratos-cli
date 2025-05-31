package cfgexp

import (
	"errors"

	"github.com/tx7do/kratos-cli/config-exporter/internal"
	"github.com/tx7do/kratos-cli/config-exporter/internal/consul"
	"github.com/tx7do/kratos-cli/config-exporter/internal/etcd"
)

func NewExporter(typeName string, endpoint string, prefix string, projectRootPath string) internal.Exporter {
	opts := &internal.Options{
		Service:     internal.ImporterType(typeName),
		Endpoint:    endpoint,
		ProjectName: prefix,
		ProjectRoot: projectRootPath,
	}

	switch opts.Service {
	default:
		fallthrough
	case internal.Consul:
		return consul.NewExporter(opts)

	case internal.Etcd:
		return etcd.NewExporter(opts)
	}
}

func Export(typeName string, endpoint string, prefix string, projectRootPath string) error {
	exporter := NewExporter(typeName, endpoint, prefix, projectRootPath)
	if exporter == nil {
		return errors.New("exporter is nil")
	}
	return exporter.Export()
}
