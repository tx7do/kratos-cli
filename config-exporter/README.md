# Local Config Files → Remote Config System Exporter

Support Remote Config Systems Exporter for Kratos CLI:

- Consul
- Etcd

## HOW TO INSTALL

```shell
go install github.com/tx7do/kratos-cli/config-exporter/cmd/cfgexp@latest
```

## HOW TO USE

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

## EXAMPLES

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
