package templates

import _ "embed"

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

//go:embed data.tpl
var DataTemplate []byte

//go:embed ent_client.tpl
var EntClientTemplate []byte

//go:embed gorm_client.tpl
var GormClientTemplate []byte

//go:embed gorm_init.tpl
var GormInitTemplate []byte

//go:embed redis_client.tpl
var RedisClientTemplate []byte
