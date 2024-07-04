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
	Client *redis.Client
	Ctx    context.Context
}

func NewRedisStorage() *RedisStorage {
	Client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	return &RedisStorage{Client: Client, Ctx: context.Background()}
}

func (rs *RedisStorage) Cleanup(keys ...string) {
	for _, key := range keys {
		rs.Client.Del(rs.Ctx, fmt.Sprintf("rate_limit:%s", key))
		rs.Client.Del(rs.Ctx, fmt.Sprintf("rate_limit:%s:blocked", key))
	}
}

func (rs *RedisStorage) IsBlocked(key string) bool {
	blockedKey := fmt.Sprintf("%s:blocked", key)
	_, err := rs.Client.Get(rs.Ctx, blockedKey).Result()
	return err != redis.Nil
}

func (rs *RedisStorage) Block(key string, blockTime time.Duration) {
	blockedKey := fmt.Sprintf("%s:blocked", key)
	rs.Client.Set(rs.Ctx, blockedKey, "1", blockTime)
}

func (rs *RedisStorage) Allow(key string, limit int, duration time.Duration) AllowResponse {
	now := time.Now()
	key = fmt.Sprintf("rate_limit:%s", key)
	blockedKey := fmt.Sprintf("%s:blocked", key)

	if rs.IsBlocked(key) {
		until, _ := rs.Client.TTL(rs.Ctx, blockedKey).Result()
		return AllowResponse{
			Allowed:     false,
			IsBlocked:   true,
			UnblockTime: now.Add(until),
		}
	}

	val, err := rs.Client.Get(rs.Ctx, key).Result()
	if err == redis.Nil {
		rs.Client.Set(rs.Ctx, key, "1", duration)
		return AllowResponse{Allowed: true}
	}
	if err != nil {
		fmt.Printf("Redis error: %v\n", err)
		return AllowResponse{Allowed: false}
	}
	count, _ := strconv.Atoi(val)
	if count < limit {
		rs.Client.Incr(rs.Ctx, key)
		return AllowResponse{Allowed: true}
	} else {
		blockTime, _ := time.ParseDuration(os.Getenv("BLOCK_TIME") + "s")
		rs.Block(key, blockTime)
		return AllowResponse{
			Allowed:     false,
			IsBlocked:   true,
			UnblockTime: now.Add(blockTime),
		}
	}
}
