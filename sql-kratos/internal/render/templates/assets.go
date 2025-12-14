package templates

import _ "embed"

//go:embed ent_repo.tpl
var EntRepoTemplate []byte

//go:embed gorm_repo.tpl
var GormRepoTemplate []byte

//go:embed service.tpl
var ServiceTemplate []byte

//go:embed main.tpl
var MainTemplate []byte

//go:embed wire.tpl
var WireTemplate []byte

//go:embed init.tpl
var InitTemplate []byte

//go:embed grpc_server.tpl
var GrpcServerTemplate []byte

//go:embed rest_server.tpl
var RestServerTemplate []byte
