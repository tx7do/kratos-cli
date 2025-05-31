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
  -v, --drv string               Database driver name to use (mysql, postgres, sqlite...) (default "mysql")
  -n, --dsn string               Data source name (connection information), for example:
                                 "mysql://user:pass@tcp(localhost:3306)/dbname"
                                 "postgres://user:pass@host:port/dbname"
  -e, --exclude-tables strings   comma-separated list of tables to exclude
  -h, --help                     help for sql2proto
  -s, --proto-path string        output path for protobuf schema files (default "./api/protos/")
  -t, --tables strings           comma-separated list of tables to inspect (all if empty)
```

## Example

```shell
sql2proto \
  --drv "postgres" \
  --dsn "postgres://postgres:pass@localhost:5432/test?sslmode=disable" \
  --proto-path "./api/protos"
```
