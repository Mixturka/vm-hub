package session

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"vm-hub/internal/config"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client  *redis.Client
	options *Options
}

func NewRedisStore(client *redis.Client, config *config.Config) *RedisStore {
	options := &Options{
		MaxAge:          config.SessionLifeTime,
		SessionName:     config.SessionName,
		SessionDomain:   config.SessionDomain,
		SessionSecure:   config.SessionSecure,
		SessionHttpOnly: config.SessionHttpOnly,
		SessionFolder:   config.SessionFolder,
		SessionSecret:   config.SessionSecret,
		CookiesSecret:   config.CookiesSecret,
	}

	return &RedisStore{
		client:  client,
		options: options,
	}
}

type Options struct {
	MaxAge          int
	SessionName     string
	SessionDomain   string
	SessionSecure   bool
	SessionHttpOnly bool
	SessionFolder   string
	SessionSecret   string
	CookiesSecret   string
}

func (rs *RedisStore) Get(ctx context.Context, r *http.Request, name string) (map[string]interface{}, error) {
	data, err := rs.client.Get(ctx, name).Result()
	if err == redis.Nil {
		return map[string]interface{}{}, nil
	} else if err != nil {
		return nil, err
	}

	var values map[string]interface{}
	err = json.Unmarshal([]byte(data), &values)
	if err != nil {
		return nil, errors.New("failed to json unmarshal session data")
	}

	return values, nil
}

func (rs *RedisStore) Save(ctx context.Context, r *http.Request, w http.ResponseWriter, name string, values map[string]interface{}) error {
	data, err := json.Marshal(values)
	if err != nil {
		return errors.New("failed to json marshal session data")
	}

	err = rs.client.Set(ctx, name, data, time.Duration(rs.options.MaxAge)).Err()
	if err != nil {
		return errors.New("failed to save session to Redis")
	}

	return nil
}

func (rs *RedisStore) Delete(ctx context.Context, r *http.Request, w http.ResponseWriter, name string) error {
	err := rs.client.Del(ctx, name).Err()
	if err != nil {
		return errors.New("failed to delete session from Redis")
	}

	return err
}
