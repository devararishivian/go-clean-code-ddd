package memory

import (
	"context"
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

var ctx = context.Background()

func (c *CacheRepositoryImpl) Set(cache entity.Cache) error {
	return c.redis.Client.Set(ctx, cache.Key, cache.Value, cache.TTL).Err()
}

func (c *CacheRepositoryImpl) HSet(cache entity.Cache) error {
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

func (c *CacheRepositoryImpl) HGet(key string) (entity.Cache, error) {
	return entity.Cache{}, nil
}
