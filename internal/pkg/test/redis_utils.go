package test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Mixturka/vm-hub/pkg/putils"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

type RedisTestUtil struct {
	t        *testing.T
	client   *redis.Client
	prefix   string
	clientMu sync.Mutex
}

func NewRedisTestUtil(t *testing.T) *RedisTestUtil {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	projRoot, err := putils.GetProjectRoot(cwd)
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}

	if err := godotenv.Load(filepath.Join(projRoot, ".env.test")); err != nil {
		log.Fatalf("Error loading environment file: %v", err)
	}

	connStr := MustGetEnv("TEST_REDIS_URL")
	prefix := generateUniquePrefix(t)
	options, err := redis.ParseURL(connStr)
	require.NoError(t, err)

	client := redis.NewClient(options)
	err = client.Ping(context.Background()).Err()
	require.NoError(t, err)

	t.Cleanup(func() {
		ctx := context.Background()
		iter := client.Scan(ctx, 0, prefix+"*", 0).Iterator()
		for iter.Next(ctx) {
			client.Del(ctx, iter.Val())
		}
		require.NoError(t, iter.Err())

		client.Close()
	})

	return &RedisTestUtil{
		t:      t,
		client: client,
		prefix: prefix,
	}
}

func (r *RedisTestUtil) Client() *redis.Client {
	r.clientMu.Lock()
	defer r.clientMu.Unlock()

	return r.client
}

func (r *RedisTestUtil) PrefixedKey(key string) string {
	return fmt.Sprintf("%s:%s", r.prefix, key)
}

func generateUniquePrefix(t *testing.T) string {
	uid := generateUniqueID(t)
	testName := strings.ToLower(strings.ReplaceAll(t.Name(), "/", "-"))
	return fmt.Sprintf("%s-%s", testName, uid)
}

func generateUniqueID(t *testing.T) string {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return base64.RawURLEncoding.EncodeToString(b)
}
