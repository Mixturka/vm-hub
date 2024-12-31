package session

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"vm-hub/internal/config"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/sessions"
)

type RedisSessionManager struct {
	client  *redis.Client
	session *sessions.CookieStore
	options *config.SessionOptions
}

func NewRedisSessionManager(client *redis.Client, cfg *config.Config) *RedisSessionManager {
	options := &config.SessionOptions{
		MaxAge:          cfg.SessionLifeTime,
		SessionName:     cfg.SessionName,
		SessionDomain:   cfg.SessionDomain,
		SessionSecure:   cfg.SessionSecure,
		SessionHttpOnly: cfg.SessionHttpOnly,
		SessionFolder:   cfg.SessionFolder,
		SessionSecret:   cfg.SessionSecret,
		CookiesSecret:   cfg.CookiesSecret,
	}

	return &RedisSessionManager{
		client:  client,
		session: sessions.NewCookieStore([]byte(cfg.SessionSecret)),
		options: options,
	}
}

func (rsm *RedisSessionManager) Get(ctx context.Context, r *http.Request, name string) (map[string]interface{}, error) {
	data, err := rsm.client.Get(ctx, name).Result()
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

func (rsm *RedisSessionManager) Save(ctx context.Context, r *http.Request, w http.ResponseWriter, name string, values map[string]interface{}) error {
	data, err := json.Marshal(values)
	if err != nil {
		return errors.New("failed to json marshal session data")
	}

	err = rsm.client.Set(ctx, name, data, time.Duration(rsm.options.MaxAge)).Err()
	if err != nil {
		return errors.New("failed to save session to Redis")
	}

	return nil
}

func (rsm *RedisSessionManager) Delete(ctx context.Context, r *http.Request, w http.ResponseWriter, name string) error {
	err := rsm.client.Del(ctx, name).Err()
	if err != nil {
		return errors.New("failed to delete session from Redis")
	}

	return err
}

func (rsm *RedisSessionManager) GetOptions() *config.SessionOptions {
	return rsm.options
}
