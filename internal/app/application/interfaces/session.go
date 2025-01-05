package interfaces

import (
	"context"
)

type SessionStorage interface {
	Get(ctx context.Context, key string) (map[string]interface{}, error)
	Set(ctx context.Context, key string, value map[string]interface{}, ttlSeconds int) error
	Delete(ctx context.Context, key string) error
}
