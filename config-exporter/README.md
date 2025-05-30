# 本地配置文件导出到远程配置中心的导出工具

支持的远程配置系统：

- Consul
- Etcd

## 使用方法

```shell
go install github.com/tx7do/kratos-cli/config-exporter

config-exporter \
    --type=etcd \
    --addr=localhost:2379 \
    --proj=kratos_admin
```

或者

```shell
go run -mod=mod github.com/tx7do/kratos-cli/config-exporter \
    --type=etcd \
    --addr=localhost:2379 \
    --proj=kratos_admin
```
