# GoWind Toolkit

**A comprehensive all-in-one toolkit for Go-Kratos microservice development, including scaffolding, CRUD code generation,
dev tools, operation utilities, CLI and desktop UI.**

一个为 Go-Kratos 微服务生态打造的**一站式全能工具集**，包含脚手架、自动化代码生成、开发辅助、运维工具、命令行与可视化桌面客户端。

| CMD                                   | Description                                                                                  |
|---------------------------------------|----------------------------------------------------------------------------------------------|
| [gowind](./gowind/README.md)          | Main CLI & UI entry: project creation, scaffolding, management, and visual operation panel   |
| [cfgexp](./config-exporter/README.md) | Export local configs to `Consul` / `Etcd` / `Nacos` for distributed configuration management |
| [sql2orm](./sql-orm/README.md)        | Generate ORM models (`ent` / `gorm`) from SQL database schemas                               |
| [sql2proto](./sql-proto/README.md)    | Generate Protobuf files (`gRPC` / `REST`) from SQL tables.                                   |
| [sql2kratos](./sql-kratos/README.md)  | One-click generation of complete Kratos microservice projects from SQL.                      |

## Tool Details

## gowind – Main Toolkit (CLI + UI)

The core entry of GoWind Toolkit.

Features:

- Project scaffolding & initialization
- Visual UI panel for development & operation
- One-click CRUD generation
- Unified command management
- Built-in tools for dev & ops
- 
#### Install

```shell
go install github.com/tx7do/go-wind-toolkit/gowind/cmd/gow@latest
```
`
### cfgexp – Configuration Exporter

Export local configuration files to remote configuration centers (Consul, Etcd, Nacos) for distributed systems.

#### Install

```shell
go install github.com/tx7do/go-wind-toolkit/config-exporter/cmd/cfgexp@latest
```

### sql2orm – SQL → ORM Generator

Import SQL schemas and generate type-safe ORM code for Kratos.

Supports:

- ent
- gorm

#### Install

```shell
go install github.com/tx7do/go-wind-toolkit/sql-orm/cmd/sql2orm@latest
```

### sql2proto – SQL → Protobuf Generator

Automatically generate Protobuf, gRPC & REST service definitions from SQL tables.

#### Install

```shell
go install github.com/tx7do/go-wind-toolkit/sql-proto/cmd/sql2proto@latest
```


## sql2kratos – SQL → Full Kratos Service

One-click generation of complete microservice project from SQL database:

- ORM models
- Protobuf
- Service layer
- Data access layer
- gRPC / REST servers

#### Install

```shell
go install github.com/tx7do/go-wind-toolkit/sql-kratos/cmd/sql2kratos@latest
```
