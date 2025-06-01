syntax = "proto3";

package {{.Package}};

import "gnostic/openapi/v3/annotations.proto";

import "google/protobuf/empty.proto";
import "google/protobuf/field_mask.proto";
import "google/protobuf/timestamp.proto";

import "pagination/v1/pagination.proto";

// {{.Comment}}服务
service {{.PascalName}}Service {
  // 查询列表
  rpc List (pagination.PagingRequest) returns (List{{.PascalName}}Response) {}

  // 查询详情
  rpc Get (Get{{.PascalName}}Request) returns ({{.PascalName}}) {}

  // 创建
  rpc Create (Create{{.PascalName}}Request) returns (google.protobuf.Empty) {}

  // 更新
  rpc Update (Update{{.PascalName}}Request) returns (google.protobuf.Empty) {}

  // 删除
  rpc Delete (Delete{{.PascalName}}Request) returns (google.protobuf.Empty) {}

  // 批量创建
  rpc BatchCreate (BatchCreate{{.PascalName}}Request) returns (BatchCreate{{.PascalName}}Response) {}
}

// {{.Comment}}
message {{.PascalName}} {
{{range .Fields}}  optional {{.Type}} {{.SnakeName}} = {{.Number}} [
    json_name = "{{.CamelName}}",
    (gnostic.openapi.v3.property) = {description: "{{.Comment}}"}
  ]; // {{.Comment}}

{{end -}}}

message List{{.PascalName}}Response {
  repeated {{.PascalName}} items = 1;
  uint32 total = 2;
}

message Get{{.PascalName}}Request {
  uint32 id = 1;
}

message Create{{.PascalName}}Request {
  {{.PascalName}} data = 1;
}

message Update{{.PascalName}}Request {
  {{.PascalName}} data = 1;

  google.protobuf.FieldMask update_mask = 2 [
    (gnostic.openapi.v3.property) = {
      description: "要更新的字段列表",
      example: {yaml : "id,realname,username"}
    },
    json_name = "updateMask"
  ]; // 要更新的字段列表

  optional bool allow_missing = 3 [
    (gnostic.openapi.v3.property) = {description: "如果设置为true的时候，资源不存在则会新增(插入)，并且在这种情况下`updateMask`字段将会被忽略。"},
    json_name = "allowMissing"
  ]; // 如果设置为true的时候，资源不存在则会新增(插入)，并且在这种情况下`updateMask`字段将会被忽略。
}

message Delete{{.PascalName}}Request {
  uint32 id = 1;
}

message BatchCreate{{.PascalName}}Request {
  repeated {{.PascalName}} data = 1;
}
message BatchCreate{{.PascalName}}Response {
  repeated {{.PascalName}} data = 1;
}
