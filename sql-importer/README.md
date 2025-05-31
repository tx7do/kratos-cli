# SQL → Schema Importer

Copy from <https://github.com/zeevmoney/entimport>

## HOW TO INSTALL

```shell
go install github.com/tx7do/kratos-cli/sql-importer/cmd/sql2orm@latest
```

## HOW TO USE

```shell
SQL to ORM code Importer is a tool to generate ORM code from SQL database schemas.

Usage:
  sql2orm [flags]

Flags:
      --dao-path string          output path for DAO code (for gorm) (default "./daos/")
      --dsn string               data source name (connection information), for example:
                                 "mysql://user:pass@tcp(localhost:3306)/dbname"
                                 "postgres://user:pass@host:port/dbname"
      --exclude-tables strings   comma-separated list of tables to exclude
  -h, --help                     help for sql2orm
      --orm string               ORM type to use (ent, gorm), default is 'ent' (default "ent")
      --schema-path string       output path for schema (default "./ent/schema/")
      --tables strings           comma-separated list of tables to inspect (all if empty)
```

## Example

for `ent` ORM:

```shell
sql2orm \
  -orm "ent" \
  -dsn "postgres://postgres:pass@localhost:5432/test?sslmode=disable" \
  -schema-path "./ent/schema"
```

for `gorm` ORM:

```shell
sql2orm \
  -orm "gorm" \
  -drv "postgres" \
  -dsn "postgres://postgres:pass@localhost:5432/test?sslmode=disable" \
  -schema-path "./daos/models" \
  -dao-path "./daos/"
```
