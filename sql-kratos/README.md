# SQL → Kratos Microservice

This tool imports the SQL database schemas and generates Kratos microservice code.

## HOW TO INSTALL

```shell
go install github.com/tx7do/kratos-cli/sql-kratos/cmd/sql2kratos@latest
```

## HOW TO USE

```shell
sql2kratos imports the SQL database schemas and generates Kratos microservice code.

Usage:
  sql2kratos [flags]

Flags:
  -n, --dsn string          Data source name (connection information), for example:
                            "mysql://user:pass@tcp(localhost:3306)/dbname"
                            "postgres://user:pass@host:port/dbname"
  -e, --excludes strings    comma-separated list of tables to exclude
  -l, --gen-data            enable generate data package code (default true)
  -k, --gen-main            enable generate main package code (default true)
  -z, --gen-orm             enable generate ORM code (default true)
  -q, --gen-proto           enable generate protobuf schema files (default true)
  -w, --gen-srv             enable generate server package code (default true)
  -a, --gen-svc             enable generate service package code (default true)
  -h, --help                help for sql2kratos
  -i, --includes strings    comma-separated list of tables to inspect (all if empty)
  -m, --module string       Target module name for the generated code, e.g., 'admin' (default "admin")
  -r, --orm string          ORM type to use (ent, gorm) (default "ent")
  -o, --output string       output path for protobuf schema files (default "./api/protos/")
  -p, --project string      Project name for the generated code, e.g., 'kratos-admin' (default "kratos-admin")
  -x, --repo                use repository pattern (default true)
  -g, --servers strings     comma-separated list of servers to generate, e.g., "grpc,rest" (default [grpc])
  -c, --service string      Service name for the generated code, e.g., 'user' (default "user")
  -s, --src-module string   Source module name, for REST service generate, e.g., "admin" (default "user")
  -v, --version string      Version of the module, e.g., 'v1' (default "v1")
```

## EXAMPLES

generate code for gRPC service:

```shell
sql2kratos \
  -p "kratos-admin" \
  -n "postgres://postgres:pass@localhost:5432/test?sslmode=disable" \
  -r "ent" \
  -o "." \
  -m "user" \
  -c "user" \
  -g "grpc"
```

generate code for REST service:

```shell
sql2kratos \
  -p "kratos-admin" \
  -n "postgres://postgres:pass@localhost:5432/test?sslmode=disable" \
  -r "ent" \
  -o "." \
  -s "user" \
  -m "admin" \
  -c "admin" \
  -g "rest" \
  -x=false \
  -z=false \
  -l=false
```
