package memory

import (
	"context"
	"errors"
	"github.com/devararishivian/go-clean-code-ddd/internal/domain/entity"
	"github.com/devararishivian/go-clean-code-ddd/internal/domain/repository"
	"github.com/devararishivian/go-clean-code-ddd/internal/infrastructure"
	"github.com/redis/go-redis/v9"
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
	value, ok := cache.Value.(map[string]interface{})
	if !ok {
		return ErrorInvalidCacheValueTypeHash
	}

	_, err := c.redis.Client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, cache.Key, value)
		pipe.Expire(ctx, cache.Key, cache.TTL)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
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

func (c *CacheRepositoryImpl) HGet(key, field string) (entity.Cache, error) {
	result := entity.Cache{
		Key: key,
	}

	val, err := c.redis.Client.HGet(ctx, key, field).Result()
	if err != nil {
		return result, err
	}

	result.Value = val

	return result, nil
}

func (c *CacheRepositoryImpl) HGetAll(key string) (entity.Cache, error) {
	result := entity.Cache{
		Key: key,
	}

	values, err := c.redis.Client.HGetAll(ctx, key).Result()
	if err != nil {
		return result, err
	}

	valMap := make(map[string]any, len(values))
	for k, v := range values {
		valMap[k] = v
	}

	result.Value = valMap

	return result, nil
}

func (c *CacheRepositoryImpl) Del(key string) error {
	return c.redis.Client.Del(ctx, key).Err()
}
