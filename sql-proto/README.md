# SQL → Protobuf

## HOW TO INSTALL

```shell
go install github.com/tx7do/kratos-cli/sql-proto/cmd/sql2proto@latest
```

## HOW TO USE

```shell
sql2proto is a tool to import SQL database schema and generate Protobuf code.

Usage:
  sql2proto [flags]

Flags:
  -n, --dsn string          Data source name (connection information), for example:
                            "mysql://user:pass@tcp(localhost:3306)/dbname"
                            "postgres://user:pass@host:port/dbname"
  -e, --excludes strings    comma-separated list of tables to exclude
  -h, --help                help for sql2proto
  -i, --includes strings    comma-separated list of tables to inspect (all if empty)
  -m, --module string       module name for the generated code, e.g., 'admin' (default "admin")
  -o, --output string       output path for protobuf schema files (default "./api/protos/")
  -s, --src-module string   Source module name, for REST service generate, e.g., "admin" (default "user")
  -t, --type string         generate RPC service type, "rest" for REST service, "grpc" for gRPC service (default "grpc")
  -v, --version string      Version of the module, e.g., 'v1' (default "v1")
```

## EXAMPLES

generate gRPC service from PostgreSQL database schema:

```shell
sql2proto \
  -n "postgres://postgres:pass@localhost:5432/test?sslmode=disable" \
  -o "./api/protos" \
  -t "grpc" \
  -m "user"
```

generate REST service from MySQL database schema:

```shell
sql2proto \
  -n "mysql://root:pass@localhost:3306/test" \
  -o "./api/protos" \
  -t "rest" \
  -m "admin" \
  -s "user"
```
