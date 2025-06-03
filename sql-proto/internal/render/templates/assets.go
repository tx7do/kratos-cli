package templates

import _ "embed"

//go:embed grpc_proto.tpl
var GrpcProtoTemplateData []byte

//go:embed rest_proto.tpl
var RestProtoTemplateData []byte
