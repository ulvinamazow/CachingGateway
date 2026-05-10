package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ulvinamazow/caching-gateway/internal/model"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

func (redisCache *RedisCache) Get(key string) (model.CachedItem, bool, error) {
	ctx, cancel := context.WithTimeout(redisCache.ctx, time.Second*2)
	defer cancel()

	value, err := redisCache.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return model.CachedItem{}, false, nil
	}

	if err != nil {
		return model.CachedItem{}, false, err
	}

	var item model.CachedItem

	err = json.Unmarshal([]byte(value), &item)
	if err != nil {
		return model.CachedItem{}, false, err
	}

	return item, true, nil
}

func (redisCache *RedisCache) Set(key string, value model.CachedItem) error {
	ctx, cancel := context.WithTimeout(redisCache.ctx, time.Second*2)
	defer cancel()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return redisCache.client.Set(ctx, key, data, time.Minute*30).Err()
}

func (redisCache *RedisCache) Clear() error {
	ctx, cancel := context.WithTimeout(redisCache.ctx, time.Second*5)
	defer cancel()

	return redisCache.client.FlushDB(ctx).Err()
}
