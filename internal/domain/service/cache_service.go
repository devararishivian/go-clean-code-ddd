package service

import "github.com/devararishivian/antrekuy/internal/domain/entity"

type CacheService interface {
	Set(cache entity.Cache) error
	HSet(cache entity.Cache) error
	Get(key string) (entity.Cache, error)
	HGet(key, field string) (entity.Cache, error)
	HGetAll(key string) (entity.Cache, error)
	Del(key string) error
}
