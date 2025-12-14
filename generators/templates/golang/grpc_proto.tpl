syntax = "proto3";

package {{.Package}};

import "gnostic/openapi/v3/annotations.proto";

import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";
import "google/protobuf/field_mask.proto";

import "pagination/v1/pagination.proto";

// {{.Comment}}服务
service {{.PascalName}}Service {
  // 查询列表
  rpc List (pagination.PagingRequest) returns (List{{.PascalName}}Response) {}

  // 查询详情
  rpc Get (Get{{.PascalName}}Request) returns ({{.PascalName}}) {}

  // 创建
  rpc Create (Create{{.PascalName}}Request) returns ({{.PascalName}}) {}

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
  uint64 total = 2;
}

message Get{{.PascalName}}Request {
  oneof query_by {
    uint32 id = 1 [
      (gnostic.openapi.v3.property) = {description: "ID", read_only: true},
      json_name = "id"
    ]; // ID
  }

  optional google.protobuf.FieldMask view_mask = 100 [
    json_name = "viewMask",
    (gnostic.openapi.v3.property) = {
      description: "视图字段过滤器，用于控制返回的字段"
    }
  ]; // 视图字段过滤器，用于控制返回的字段
}

message Create{{.PascalName}}Request {
  {{.PascalName}} data = 1;
}

message Update{{.PascalName}}Request {
  uint32 id = 1;

  {{.PascalName}} data = 2;

  google.protobuf.FieldMask update_mask = 100 [
    (gnostic.openapi.v3.property) = {
      description: "要更新的字段列表",
      example: {yaml : "id,realname,username"}
    },
    json_name = "updateMask"
  ]; // 要更新的字段列表

  optional bool allow_missing = 101 [
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
