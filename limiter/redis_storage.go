package limiter

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage() *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	return &RedisStorage{client: client, ctx: context.Background()}
}

func (rs *RedisStorage) Allow(key string, limit int, duration time.Duration) bool {
	key = fmt.Sprintf("rate_limit:%s", key)
	val, err := rs.client.Get(rs.ctx, key).Result()
	if err == redis.Nil {
		rs.client.Set(rs.ctx, key, "1", duration)
		return true
	}
	if err != nil {
		fmt.Printf("Redis error: %v\n", err)
		return false
	}
	count, _ := strconv.Atoi(val)
	if count < limit {
		rs.client.Incr(rs.ctx, key)
		return true
	}
	return false
}
