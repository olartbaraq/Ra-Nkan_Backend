package redisInit

import (
	"log"

	"github.com/olartbaraq/spectrumshelf/utils"
	"github.com/redis/go-redis/v9"
)

func RedisInit() *redis.Client {

	otherConfig, err := utils.LoadOtherConfig(".")
	if err != nil {
		log.Fatal("Could not load env config", err)
	}
	var Rdb = redis.NewClient(&redis.Options{
		Addr:     otherConfig.RedisAddress,
		Password: otherConfig.RedisPassword,
		DB:       0, // use default DB
	})

	return Rdb

}
