package interfaces

import (
	"context"
	"net/http"
)

type SessionManager interface {
	Get(ctx context.Context, r *http.Request, name string) (map[string]interface{}, error)
	Save(ctx context.Context, r *http.Request, w http.ResponseWriter, name string, values map[string]interface{}) error
	Delete(ctx context.Context, r *http.Request, w http.ResponseWriter, name string) error
}
