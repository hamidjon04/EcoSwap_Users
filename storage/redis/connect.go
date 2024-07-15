package redis

import (
	"ecoswap/config"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis()*redis.Client{
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Load().REDIS_PORT,
		Password: "",
		DB: 0,
	})

	return rdb
}



