package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	redisClient "github.com/tx7do/kratos-bootstrap/cache/redis"
)

// NewRedisClient 创建Redis客户端
func NewRedisClient(cfg *conf.Bootstrap, logger log.Logger) *redis.Client {
	l := log.NewHelper(log.With(logger, "module", "redis/data/{{.Service}}-service"))
	return redisClient.NewClient(cfg.Data, l)
}
