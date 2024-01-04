package redisInit

import (
	config "github.com/olartbaraq/spectrumshelf/configs"
	"github.com/redis/go-redis/v9"
)

func RedisInit() *redis.Client {
	var Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: config.EnvRedisPassword(),
		DB:       0, // use default DB
	})

	return Rdb

}
