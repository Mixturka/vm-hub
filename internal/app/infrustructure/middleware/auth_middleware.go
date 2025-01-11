package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Mixturka/vm-hub/internal/app/application/services"
	"github.com/Mixturka/vm-hub/internal/app/infrustructure/session"
)

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware(userService *services.UserService, sessionManager *session.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values, err := sessionManager.GetSession(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, ok := values["userID"]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := userService.FindByID(ctx, userID.(string))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
