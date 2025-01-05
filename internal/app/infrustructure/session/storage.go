package session

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

func (rs *RedisStore) Get(ctx context.Context, key string) (map[string]interface{}, error) {
	data, err := rs.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var values map[string]interface{}
	if err := json.Unmarshal([]byte(data), &values); err != nil {
		return nil, errors.New("failed to unmarshal json session data")
	}

	return values, nil
}

func (rs *RedisStore) Set(ctx context.Context, key string, value map[string]interface{}, ttlSeconds int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.New("failed to marshal json session data")
	}

	return rs.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err()
}

func (rs *RedisStore) Delete(ctx context.Context, key string) error {
	return rs.client.Del(ctx, key).Err()
}
