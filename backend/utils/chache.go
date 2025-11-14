// file: utils/cache.go
package utils

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	RedisCtx    = context.Background()
)

func InitRedis(addr, password string, db int) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	println("Redis connected in" + addr)
	return RedisClient.Ping(RedisCtx).Err()
}

func GetOTPFromCache(phone string) (string, bool, error) {
	if RedisClient == nil {
		return "", false, nil
	}
	val, err := RedisClient.Get(RedisCtx, "otp:"+phone).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func SetOTPToCache(phone, code string, ttl time.Duration) error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Set(RedisCtx, "otp:"+phone, code, ttl).Err()
}

func DeleteOTPFromCache(phone string) error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Del(RedisCtx, "otp:"+phone).Err()
}
