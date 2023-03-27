package repository

import "github.com/devararishivian/go-clean-code-ddd/internal/domain/entity"

type CacheRepository interface {
	Set(cache entity.Cache) error
	HSet(cache entity.Cache) error
	Get(key string) (entity.Cache, error)
	HGet(key, field string) (entity.Cache, error)
	HGetAll(key string) (entity.Cache, error)
	Del(key string) error
}
