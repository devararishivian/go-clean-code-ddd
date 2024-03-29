package memory_test

import (
	"errors"
	"github.com/go-redis/redismock/v9"
	"testing"
	"time"

	"github.com/devararishivian/go-clean-code-ddd/internal/domain/entity"
	"github.com/devararishivian/go-clean-code-ddd/internal/infrastructure"
	"github.com/devararishivian/go-clean-code-ddd/internal/infrastructure/memory"
	"github.com/stretchr/testify/assert"
)

func TestCacheRepositoryImpl(t *testing.T) {
	redisClient, mock := redismock.NewClientMock()
	redis := &infrastructure.Redis{
		Client: redisClient,
	}

	cacheRepository := memory.NewCacheRepository(redis)

	t.Run("Set", func(t *testing.T) {
		cache := entity.Cache{
			Key:   "key",
			Value: "value",
			TTL:   10 * time.Second,
		}

		mock.ExpectSet(cache.Key, cache.Value, cache.TTL).SetVal("OK")

		err := cacheRepository.Set(cache)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HSet", func(t *testing.T) {
		cache := entity.Cache{
			Key: "key",
			Value: map[string]interface{}{
				"field1": "value1",
				"field2": 2,
				"field3": true,
			},
			TTL: 10 * time.Second,
		}

		mock.ExpectTxPipeline()
		mock.ExpectHSet(cache.Key, cache.Value).SetVal(1)
		mock.ExpectExpire(cache.Key, cache.TTL).SetVal(true)
		mock.ExpectTxPipelineExec()

		err := cacheRepository.HSet(cache)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HSet invalid value type", func(t *testing.T) {
		cache := entity.Cache{
			Key:   "key",
			Value: "invalid",
			TTL:   10 * time.Second,
		}

		err := cacheRepository.HSet(cache)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, memory.ErrorInvalidCacheValueTypeHash))
	})

	t.Run("Get", func(t *testing.T) {
		key := "key"
		value := "value"

		mock.ExpectGet(key).SetVal(value)

		result, err := cacheRepository.Get(key)

		assert.NoError(t, err)
		assert.Equal(t, key, result.Key)
		assert.Equal(t, value, result.Value)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HGet", func(t *testing.T) {
		key := "key"
		field := "field"
		value := "value"

		mock.ExpectHGet(key, field).SetVal("value")

		result, err := cacheRepository.HGet(key, field)

		assert.NoError(t, err)
		assert.Equal(t, key, result.Key)
		assert.Equal(t, value, result.Value)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HGetAll", func(t *testing.T) {
		key := "key"
		value := map[string]interface{}{
			"field1": "value1",
			"field2": "2",
			"field3": "true",
		}

		mock.ExpectHGetAll(key).SetVal(map[string]string{
			"field1": "value1",
			"field2": "2",
			"field3": "true",
		})

		result, err := cacheRepository.HGetAll(key)

		assert.NoError(t, err)
		assert.Equal(t, key, result.Key)
		assert.Equal(t, value, result.Value)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
