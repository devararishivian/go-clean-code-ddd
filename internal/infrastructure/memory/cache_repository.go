package memory

import (
	"context"
	"errors"
	"github.com/devararishivian/antrekuy/internal/domain/entity"
	"github.com/devararishivian/antrekuy/internal/domain/repository"
	"github.com/devararishivian/antrekuy/internal/infrastructure"
)

type CacheRepositoryImpl struct {
	redis *infrastructure.Redis
}

func NewCacheRepository(redis *infrastructure.Redis) repository.CacheRepository {
	return &CacheRepositoryImpl{
		redis: redis,
	}
}

var (
	ctx                            = context.Background()
	ErrorInvalidCacheValueTypeHash = errors.New("invalid cache value type, expected map[string]interface{}")
)

func (c *CacheRepositoryImpl) Set(cache entity.Cache) error {
	return c.redis.Client.Set(ctx, cache.Key, cache.Value, cache.TTL).Err()
}

func (c *CacheRepositoryImpl) HSet(cache entity.Cache) error {
	// Convert the value to a map[string]interface{} to set as a hash field
	value, ok := cache.Value.(map[string]interface{})
	if !ok {
		return ErrorInvalidCacheValueTypeHash
	}

	return c.redis.Client.HSet(ctx, cache.Key, value).Err()
}

func (c *CacheRepositoryImpl) Get(key string) (result entity.Cache, err error) {
	val, err := c.redis.Client.Get(ctx, key).Result()
	if err != nil {
		return result, err
	}

	result.Key = key
	result.Value = val

	return result, nil
}

func (c *CacheRepositoryImpl) HGet(key string) (entity.Cache, error) {
	result := entity.Cache{
		Key: key,
	}

	// Use the Redis HGetAll command to get all the hash fields and values
	vals, err := c.redis.Client.HGetAll(ctx, key).Result()
	if err != nil {
		return result, err
	}

	// Convert the map[string]string returned by HGetAll to map[string]interface{}
	valMap := make(map[string]interface{}, len(vals))
	for k, v := range vals {
		valMap[k] = v
	}

	result.Value = valMap

	return result, nil
}
