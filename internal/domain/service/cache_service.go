package service

import "github.com/devararishivian/antrekuy/internal/domain/entity"

type CacheService interface {
	Set(cache entity.Cache) error
	HSet(cache entity.Cache) error
	Get(key string) (entity.Cache, error)
	HGet(key string) (entity.Cache, error)
}
