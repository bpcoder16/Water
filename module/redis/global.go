package redis

import "github.com/redis/go-redis/v9"

var defaultRedis *redis.Client

func GetDefaultRedis() *redis.Client {
	return defaultRedis
}
