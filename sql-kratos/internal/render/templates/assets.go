package templates

import _ "embed"

//go:embed ent_data.tpl
var EntDataTemplateData []byte

//go:embed gorm_data.tpl
var GormDataTemplateData []byte

//go:embed service.tpl
var ServiceTemplateData []byte

//go:embed main.tpl
var MainTemplateData []byte

//go:embed wire.tpl
var WireTemplateData []byte

//go:embed init.tpl
var InitTemplateData []byte

//go:embed grpc_server.tpl
var GrpcTemplateServerData []byte

//go:embed rest_server.tpl
var RestTemplateServerData []byte
