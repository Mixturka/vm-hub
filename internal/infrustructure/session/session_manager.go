package session

import (
	"context"
	"errors"
	"net/http"

	"github.com/Mixturka/vm-hub/internal/application/interfaces"
	"github.com/Mixturka/vm-hub/internal/config"

	"github.com/google/uuid"
)

type SessionManager struct {
	storage interfaces.SessionStorage
	options *config.SessionOptions
}

func NewSessionManager(storage interfaces.SessionStorage, options *config.SessionOptions) *SessionManager {
	return &SessionManager{
		storage: storage,
		options: options,
	}
}

func (sm *SessionManager) CreateSession(w http.ResponseWriter, values map[string]interface{}) (string, error) {
	sessionID := uuid.NewString()
	err := sm.storage.Set(context.Background(), sessionID, values, sm.options.MaxAge)
	if err != nil {
		return "", errors.New("failed to save session")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sm.options.SessionName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: sm.options.SessionHttpOnly,
		Secure:   sm.options.SessionSecure,
		MaxAge:   sm.options.MaxAge,
		Domain:   sm.options.SessionDomain,
	})

	return sessionID, nil
}

func (sm *SessionManager) GetSession(r *http.Request) (map[string]interface{}, error) {
	cookie, err := r.Cookie(sm.options.SessionName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
		return nil, err
	}

	return sm.storage.Get(context.Background(), cookie.Value)
}

func (sm *SessionManager) DestroySession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sm.options.SessionName)
	if err != nil {
		return nil
	}

	err = sm.storage.Delete(context.Background(), cookie.Value)
	if err != nil {
		return errors.New("failed to delete session")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sm.options.SessionName,
		Value:    "",
		Path:     "/",
		Domain:   sm.options.SessionDomain,
		HttpOnly: sm.options.SessionHttpOnly,
		Secure:   sm.options.SessionSecure,
		MaxAge:   -1,
	})

	return nil
}
