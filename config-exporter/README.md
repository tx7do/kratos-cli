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
  -a, --addr string   remote config service address (default "127.0.0.1:8500")
  -h, --help          help for cfgexp
  -p, --proj string   project name, this name is used to key prefix in remote config service
  -r, --root string   project root dir (default "./")
  -t, --type string   remote config service name (consul, etcd, etc.) (default "consul")
```

## Example

for `etcd` remote config service:

```shell
cfgexp \
    --type "etcd" \
    --addr "localhost:2379" \
    --proj "kratos_admin"
```

for `consul` remote config service:

```shell
cfgexp \
    --type "consul" \
    --addr "localhost:8500" \
    --proj "kratos_admin"
```
