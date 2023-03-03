package infrastructure

import (
	appConfig "github.com/devararishivian/antrekuy/internal/config"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis() (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     appConfig.Cache.Address,
		Password: appConfig.Cache.Password,
		DB:       appConfig.Cache.DB,
	})

	return &Redis{client}, nil
}

func (r *Redis) Close() error {
	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			return err
		}
	}

	return nil
}
