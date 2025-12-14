syntax = "proto3";

package {{.TargetPackage}};

import "gnostic/openapi/v3/annotations.proto";
import "google/protobuf/empty.proto";
import "google/api/annotations.proto";
import "pagination/v1/pagination.proto";

import "{{.SourceProto}}";

// {{.Comment}}管理服务
service {{.PascalName}}Service {
  // 查询列表
  rpc List (pagination.PagingRequest) returns ({{.SourcePackage}}.List{{.PascalName}}Response) {
    option (google.api.http) = {
        get: "{{.Path}}"
    };
  }

  // 查询详情
  rpc Get ({{.SourcePackage}}.Get{{.PascalName}}Request) returns ({{.SourcePackage}}.{{.PascalName}}) {
    option (google.api.http) = {
        get: "{{.Path}}/{id}"
    };
  }

  // 创建
  rpc Create ({{.SourcePackage}}.Create{{.PascalName}}Request) returns ({{.SourcePackage}}.{{.PascalName}}) {
    option (google.api.http) = {
        post: "{{.Path}}"
        body: "*"
    };
  }

  // 更新
  rpc Update ({{.SourcePackage}}.Update{{.PascalName}}Request) returns (google.protobuf.Empty) {
    option (google.api.http) = {
        put: "{{.Path}}/{id}"
        body: "*"
    };
  }

  // 删除
  rpc Delete ({{.SourcePackage}}.Delete{{.PascalName}}Request) returns (google.protobuf.Empty) {
    option (google.api.http) = {
        delete: "{{.Path}}/{id}"
    };
  }

  // 批量创建
  rpc BatchCreate ({{.SourcePackage}}.BatchCreate{{.PascalName}}Request) returns ({{.SourcePackage}}.BatchCreate{{.PascalName}}Response) {
    option (google.api.http) = {
        post: "{{.Path}}/batch"
        body: "*"
    };
  }
}
