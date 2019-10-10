package Storages

import (
	"github.com/go-redis/redis"
	redisStorage "github.com/khorevaa/EnchantedTBot/Storages/Redis"
)

func NewRedisStorage(opts redis.Options) *redisStorage.RedisStorage {

	return &redisStorage.RedisStorage{
		Client: redis.NewClient(&opts),
	}

}
