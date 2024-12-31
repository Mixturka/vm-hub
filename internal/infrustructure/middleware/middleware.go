package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"vm-hub/internal/application/interfaces"
)

type contextKey string

func SessionMiddleware(sm interfaces.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionData, err := sm.Get(r.Context(), r, "user-session")
		if err != nil {
			slog.Error("Error retrieving session", "error", err.Error())
			http.Error(w, "Error retrieving session", http.StatusInternalServerError)
			return
		}

		sessionKey := contextKey("sessionData")
		ctx := context.WithValue(r.Context(), sessionKey, sessionData)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		if err := sm.Save(r.Context(), r, w, "user-session", sessionData); err != nil {
			slog.Error("Error saving session", "error", err.Error())
			http.Error(w, "Error saving session", http.StatusInternalServerError)
		}
	})
}
