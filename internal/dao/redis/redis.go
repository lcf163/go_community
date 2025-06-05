package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"go_community/global"
)

var (
	client *redis.Client
	//Nil    = redis.Nil
)

// Init 初始化连接
func Init(cfg *global.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
		PoolSize: cfg.PoolSize,
	})
	_, err = client.Ping().Result()
	return
}

func Close() {
	_ = client.Close()
}
