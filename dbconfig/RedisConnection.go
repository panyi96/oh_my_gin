package dbconfig

import (
	"github.com/redis/go-redis/v9"
	"ohmygin/nacosconfig"
)

var RedisClient *redis.Client

func RedisConnect() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         nacosconfig.Config.Redis.Addr,
		Password:     nacosconfig.Config.Redis.Password, // no password set
		DB:           nacosconfig.Config.Redis.DB,       // use default DB
		MaxIdleConns: 100,
	})
}
