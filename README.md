# kratos-cli

[go-kratos](https://go-kratos.dev) is a modern web framework for building microservices in Go. This CLI tool provides a
convenient way to manage and interact with Kratos projects.

| CMD                                   | Description                                                                                                                                                        |
|---------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [cfgexp](./config-exporter/README.md) | This tool exports local configuration files to remote configuration systems like Consul or Etcd, making it easier to manage configurations in distributed systems. |
| [sql2orm](./sql-orm/README.md)        | This tool imports the SQL database schemas and generates ORM code for use in Kratos microservices, supporting both `ent` and `gorm` ORMs.                          |
| [sql2proto](./sql-proto/README.md)    | This tool imports the SQL database schemas and generates Protobuf code for use in Kratos microservices.                                                            |
| [sql2kratos](./sql-kratos/README.md)  | This tool imports the SQL database schemas and generates Kratos microservice code.                                                                                 |

## Local Config Files → Remote Config System Exporter

This tool exports local configuration files to remote configuration systems like Consul or Etcd, making it easier to
manage configurations in distributed systems.

Support Remote Config Systems Exporter for Kratos CLI:

- Consul
- Etcd

### HOW TO INSTALL

```shell
go install github.com/tx7do/kratos-cli/config-exporter/cmd/cfgexp@latest
```

### HOW TO USE

```shell
Config Exporter is a tool to export configuration from remote services like Consul or Etcd to local files.

Usage:
  cfgexp [flags]

Flags:
  -a, --addr string    remote config service address (default "127.0.0.1:8500")
  -e, --env string     environment name, like dev, test, prod, etc. (default "dev")
  -g, --group string   group name, this name is used to key prefix in remote config service (default "DEFAULT_GROUP")
  -h, --help           help for cfgexp
  -n, --ns string      namespace ID, used for Nacos (default "public")
  -p, --proj string    project name, this name is used to key prefix in remote config service
  -r, --root string    project root dir (default "./")
  -t, --type string    remote config service name (consul, etcd, etc.) (default "consul")
```

### EXAMPLES

for `etcd` remote config service:

```shell
cfgexp \
    -t "etcd" \
    -a "localhost:2379" \
    -p "kratos_admin"
```

for `consul` remote config service:

```shell
cfgexp \
    -t "consul" \
    -a "localhost:8500" \
    -p "kratos_admin"
```

for `nacos` remote config service:

```shell
cfgexp \
    -t "nacos" \
    -a "localhost:8848" \
    -p "kratos_admin" \
    -n "public" \
    -e "dev" \
    -g "DEFAULT_GROUP"
```

## SQL → ORM Schema Importer

This tool imports the SQL database schemas and generates ORM code for use in Kratos microservices, supporting both `ent`
and `gorm` ORMs.

### HOW TO INSTALL

```shell
go install github.com/tx7do/kratos-cli/sql-orm/cmd/sql2orm@latest
```

### HOW TO USE

```shell
sql2orm is a tool to generate ORM code from SQL database schemas.

Usage:
  sql2orm [flags]

Flags:
  -d, --dao-path string          output path for DAO code (for gorm) (default "./daos/")
  -v, --drv string               Database driver name to use (mysql, postgres, sqlite...) (default "mysql")
  -n, --dsn string               Data source name (connection information), for example:
                                 "mysql://user:pass@tcp(localhost:3306)/dbname"
                                 "postgres://user:pass@host:port/dbname"
  -e, --exclude-tables strings   comma-separated list of tables to exclude
  -h, --help                     help for sql2orm
  -o, --orm string               ORM type to use (ent, gorm) (default "ent")
  -s, --schema-path string       output path for schema (default "./ent/schema/")
  -t, --tables strings           comma-separated list of tables to inspect (all if empty)
```

### EXAMPLES

for `ent` ORM:

```shell
sql2orm \
  --orm "ent" \
  --dsn "postgres://postgres:pass@localhost:5432/test?sslmode=disable" \
  --schema-path "./ent/schema"
```

for `gorm` ORM:

```shell
sql2orm \
  --orm "gorm" \
  --drv "postgres" \
  --dsn "postgres://postgres:pass@localhost:5432/test?sslmode=disable" \
  --schema-path "./daos/models" \
  --dao-path "./daos/"
```

## SQL → Protobuf

This tool imports the SQL database schemas and generates Protobuf code for use in Kratos microservices.

### HOW TO INSTALL

```shell
go install github.com/tx7do/kratos-cli/sql-proto/cmd/sql2proto@latest
```

### HOW TO USE

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

### EXAMPLES

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

## SQL → Kratos Microservice

This tool imports the SQL database schemas and generates Kratos microservice code.

### HOW TO INSTALL

```shell
go install github.com/tx7do/kratos-cli/sql-kratos/cmd/sql2kratos@latest
```

### HOW TO USE

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

### EXAMPLES

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
