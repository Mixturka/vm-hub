package session_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Mixturka/vm-hub/internal/app/infrustructure/session"
	"github.com/Mixturka/vm-hub/internal/pkg/test"
	"github.com/Mixturka/vm-hub/internal/pkg/test/mock"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Unit tests
func TestRedisStore_Get_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockCmdable(ctrl)
	ctx := context.Background()
	key := "test_session"
	expectedValues := map[string]interface{}{"key": "value"}

	data, _ := json.Marshal(expectedValues)
	mockClient.EXPECT().Get(ctx, key).Return(redis.NewStringResult(string(data), nil))

	store := session.NewRedisStore(mockClient)

	values, err := store.Get(ctx, key)

	assert.NoError(t, err)
	assert.Equal(t, expectedValues, values)
}

func TestRedisStore_Get_NotFound(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockCmdable(ctrl)
	ctx := context.Background()
	key := "missing_key"

	mockClient.EXPECT().Get(ctx, key).Return(redis.NewStringResult("", redis.Nil))

	store := session.NewRedisStore(mockClient)

	values, err := store.Get(ctx, key)

	assert.NoError(t, err)
	assert.Nil(t, values)
}

func TestRedisStore_Get_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockCmdable(ctrl)
	ctx := context.Background()
	key := "test_session"

	mockClient.EXPECT().Get(ctx, key).Return(redis.NewStringResult("", errors.New("redis error")))

	store := session.NewRedisStore(mockClient)

	values, err := store.Get(ctx, key)

	assert.Error(t, err)
	assert.Nil(t, values)
}

func TestRedisStore_Set_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockCmdable(ctrl)
	ctx := context.Background()
	key := "test_session"
	ttlSeconds := 60
	values := map[string]interface{}{"key": "value"}

	data, _ := json.Marshal(values)
	mockClient.EXPECT().Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Return(redis.NewStatusResult("", nil))

	store := session.NewRedisStore(mockClient)

	err := store.Set(ctx, key, values, ttlSeconds)

	assert.NoError(t, err)
}

func TestRedisStore_Set_MarshalError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockCmdable(ctrl)
	ctx := context.Background()
	key := "test_session"
	ttlSeconds := 60
	values := map[string]interface{}{"key": make(chan int)}

	store := session.NewRedisStore(mockClient)

	err := store.Set(ctx, key, values, ttlSeconds)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal json session data")
}

func TestRedisStore_Set_RedisError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockCmdable(ctrl)
	ctx := context.Background()
	key := "test_session"
	ttlSeconds := 60
	values := map[string]interface{}{"key": "value"}

	data, _ := json.Marshal(values)
	mockClient.EXPECT().Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Return(redis.NewStatusResult("", errors.New("redis error")))

	store := session.NewRedisStore(mockClient)

	err := store.Set(ctx, key, values, ttlSeconds)

	assert.Error(t, err)
}

func TestRedisStore_Delete_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockCmdable(ctrl)
	ctx := context.Background()
	key := "test_session"

	mockClient.EXPECT().Del(ctx, key).Return(redis.NewIntResult(1, nil))

	store := session.NewRedisStore(mockClient)

	err := store.Delete(ctx, key)

	assert.NoError(t, err)
}

func TestRedisStore_Delete_RedisError(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock.NewMockCmdable(ctrl)
	ctx := context.Background()
	key := "test_session"

	mockClient.EXPECT().Del(ctx, key).Return(redis.NewIntResult(0, errors.New("redis error")))

	store := session.NewRedisStore(mockClient)

	err := store.Delete(ctx, key)

	assert.Error(t, err)
}

// Integrational tests
func TestRedisStore(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("Set and Get a key-value pair", func(t *testing.T) {
		t.Parallel()

		util := test.NewRedisTestUtil(t)
		client := util.Client()
		rs := session.NewRedisStore(client)

		key := "test-key"
		value := map[string]interface{}{
			"field1": "value1",
			"field2": float64(42),
		}

		err := rs.Set(ctx, key, value, 60)
		require.NoError(t, err)

		result, err := rs.Get(ctx, key)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, value, result)
	})

	t.Run("Get returns nil for non-existent key", func(t *testing.T) {
		t.Parallel()

		util := test.NewRedisTestUtil(t)
		client := util.Client()
		rs := session.NewRedisStore(client)

		key := "non-existent-key"
		result, err := rs.Get(ctx, key)
		require.NoError(t, err)
		require.Nil(t, result)
	})

	t.Run("Delete a key", func(t *testing.T) {
		t.Parallel()

		util := test.NewRedisTestUtil(t)
		client := util.Client()
		rs := session.NewRedisStore(client)

		key := "key-to-delete"
		value := map[string]interface{}{
			"field": "to-delete",
		}

		err := rs.Set(ctx, key, value, 60)
		require.NoError(t, err)

		err = rs.Delete(ctx, key)
		require.NoError(t, err)

		result, err := rs.Get(ctx, key)
		require.NoError(t, err)
		require.Nil(t, result)
	})

	t.Run("Set a key with expiration", func(t *testing.T) {
		t.Parallel()

		util := test.NewRedisTestUtil(t)
		client := util.Client()
		rs := session.NewRedisStore(client)

		key := "expiring-key"
		value := map[string]interface{}{
			"field": "expires",
		}

		err := rs.Set(ctx, key, value, 2)
		require.NoError(t, err)

		result, err := rs.Get(ctx, key)
		require.NoError(t, err)
		require.NotNil(t, result)

		time.Sleep(3 * time.Second)

		result, err = rs.Get(ctx, key)
		require.NoError(t, err)
		require.Nil(t, result)
	})
}
