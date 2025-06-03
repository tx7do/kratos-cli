# kratos-cli

[go-kratos](https://go-kratos.dev) is a modern web framework for building microservices in Go. This CLI tool provides a
convenient way to manage and interact with Kratos projects.

| CMD                                   | Description                                                                                                                                                        |
|---------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [cfgexp](./config-exporter/README.md) | This tool exports local configuration files to remote configuration systems like Consul or Etcd, making it easier to manage configurations in distributed systems. |
| [sql2orm](./sql-orm/README.md)        | This tool imports the SQL database schemas and generates ORM code for use in Kratos microservices, supporting both `ent` and `gorm` ORMs.                          |
| [sql2proto](./sql-proto/README.md)    | This tool imports the SQL database schemas and generates Protobuf code for use in Kratos microservices.                                                            |
| [sql2kratos](./sql-kratos/README.md)  | This tool imports the SQL database schemas and generates Kratos microservice code.                                                                                 |
