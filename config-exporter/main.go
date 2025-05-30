package main

import (
	"flag"

	"github.com/tx7do/kratos-cli/config-exporter/internal"
	"github.com/tx7do/kratos-cli/config-exporter/internal/consul"
	"github.com/tx7do/kratos-cli/config-exporter/internal/etcd"
)

func NewImporter() internal.Importer {
	var opts internal.Options

	flag.StringVar((*string)(&opts.Service), "type", "consul", "remote config service name, eg: -type consul")
	flag.StringVar(&opts.Endpoint, "addr", "127.0.0.1:8500", "consul address, eg: -addr 127.0.0.1:8500")
	flag.StringVar(&opts.ProjectName, "proj", "kratos_admin", "project name, eg: -proj kratos_admin")
	flag.StringVar(&opts.ProjectRoot, "root", "./", "project root dir, eg: -root ./")

	flag.Parse()

	switch opts.Service {
	default:
		fallthrough
	case internal.Consul:
		return consul.NewImporter(&opts)

	case internal.Etcd:
		return etcd.NewImporter(&opts)
	}
}

func main() {
	i := NewImporter()
	_ = i.Import()
}
