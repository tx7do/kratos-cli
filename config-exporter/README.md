# Local Config Files → Remote Config System Exporter

Support Remote Config Systems Exporter for Kratos CLI:

- Consul
- Etcd

## Hot To Use

```shell
# Install cfgexp
go install github.com/tx7do/kratos-cli/config-exporter/cmd/cfgexp@latest

# Export local config files to remote config system
cfgexp \
    --type=etcd \
    --addr=localhost:2379 \
    --proj=kratos_admin
```

Or

```shell
# Use go run to execute directly
go run -mod=mod github.com/tx7do/kratos-cli/config-exporter/cmd/cfgexp \
    --type=etcd \
    --addr=localhost:2379 \
    --proj=kratos_admin
```
