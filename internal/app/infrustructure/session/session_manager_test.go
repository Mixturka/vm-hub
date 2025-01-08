package session_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mixturka/vm-hub/internal/app/config"
	"github.com/Mixturka/vm-hub/internal/app/infrustructure/session"
	"github.com/Mixturka/vm-hub/internal/pkg/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateSessionSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock.NewMockSessionStorage(ctrl)
	mockStorage.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	options := &config.SessionOptions{
		SessionName:     "session_id",
		SessionHttpOnly: true,
		SessionSecure:   false,
		SessionDomain:   "localhost",
		MaxAge:          3600,
	}

	sm := session.NewSessionManager(mockStorage, options)
	w := httptest.NewRecorder()
	values := map[string]interface{}{"key": "value"}

	sessionID, err := sm.CreateSession(w, values)

	assert.NoError(t, err)
	assert.NotEmpty(t, sessionID)

	cookie := w.Result().Cookies()[0]
	assert.Equal(t, options.SessionName, cookie.Name)
	assert.Equal(t, sessionID, cookie.Value)
}

func TestCreateSessionFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock.NewMockSessionStorage(ctrl)
	mockStorage.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("storage error"))

	options := &config.SessionOptions{
		SessionName:     "session_id",
		SessionHttpOnly: true,
		SessionSecure:   false,
		SessionDomain:   "localhost",
		MaxAge:          3600,
	}

	sm := session.NewSessionManager(mockStorage, options)
	w := httptest.NewRecorder()
	values := map[string]interface{}{"key": "value"}

	sessionID, err := sm.CreateSession(w, values)

	assert.Error(t, err)
	assert.Empty(t, sessionID)
	assert.EqualError(t, err, "failed to save session")
}

func TestSessionExpiry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mock.NewMockSessionStorage(ctrl)
	mockStorage.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)

	options := &config.SessionOptions{
		SessionName:     "session_id",
		SessionHttpOnly: true,
		SessionSecure:   false,
		SessionDomain:   "localhost",
		MaxAge:          3600,
	}

	sm := session.NewSessionManager(mockStorage, options)

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  options.SessionName,
		Value: "expired_session_id",
	})

	values, err := sm.GetSession(r)

	assert.NoError(t, err)
	assert.Nil(t, values)
}
