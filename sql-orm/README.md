# SQL → ORM Schema Importer

This tool imports the SQL database schemas and generates ORM code for use in Kratos microservices, supporting both `ent` and `gorm` ORMs.

## HOW TO INSTALL

```shell
go install github.com/tx7do/kratos-cli/sql-orm/cmd/sql2orm@latest
```

## HOW TO USE

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

## EXAMPLES

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

## ACKNOWLEDGEMENT

- [ent](https://entgo.io) generator code is based on Copy from <https://github.com/zeevmoney/entimport>.
- [gorm](https://gorm.io) generator code is based on [GORM GEN](https://gorm.io/gen/index.html)
- sql parser is <https://github.com/blastrain/vitess-sqlparser>
