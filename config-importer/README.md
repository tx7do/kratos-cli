# 远程配置导入工具

支持的远程配置系统：

- Consul
- Etcd

## 使用方法

```shell
go install github.com/tx7do/kratos-cli/config-importer/cmd@latest
```

go run -mod=mod github.com/tx7do/kratos-cli/config-importer/cmd \
    --type=etcd \
    --addr=localhost:2379 \
    --proj=kratos_admin
```
