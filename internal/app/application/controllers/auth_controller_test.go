package controllers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Mixturka/vm-hub/internal/app/application/controllers"
	"github.com/Mixturka/vm-hub/internal/app/application/dtos"
	"github.com/Mixturka/vm-hub/internal/app/application/services"
	"github.com/Mixturka/vm-hub/internal/app/config"
	"github.com/Mixturka/vm-hub/internal/app/infrustructure/database/postgres"
	"github.com/Mixturka/vm-hub/internal/app/infrustructure/session"
	"github.com/Mixturka/vm-hub/internal/pkg/test"
	"github.com/Mixturka/vm-hub/pkg/putils"
	"github.com/stretchr/testify/assert"
)

var (
	migrationsPath         string
	absoluteMigrationsPath string
)

func TestMain(t *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	projRoot, err := putils.GetProjectRoot(cwd)
	if err != nil {
		log.Fatalf("Error finding project root: %v", err)
	}

	migrationsPath = test.MustGetEnv("POSTGRES_MIGRATIONS_PATH")
	absoluteMigrationsPath = test.GetAbsolutePath(projRoot, migrationsPath)

	exitCode := t.Run()
	os.Exit(exitCode)
}

// Integrational tests
func TestRegister_Success(t *testing.T) {
	t.Parallel()

	t.Run("Test register success", func(t *testing.T) {
		t.Parallel()

		ptUtil := test.NewPostgresTestUtilWithIsolatedSchema(t)
		test.ApplyMigrations(ptUtil.DB().Config().ConnString(), absoluteMigrationsPath)
		repo := postgres.NewPostgresUserRepository(ptUtil.DB())
		userService := services.NewUserService(repo)

		util := test.NewRedisTestUtil(t)
		client := util.Client()
		rs := session.NewRedisStore(client)

		sessionManager := session.NewSessionManager(rs, &config.SessionOptions{
			MaxAge:          1,
			SessionName:     "Name",
			SessionDomain:   "",
			SessionSecure:   false,
			SessionHttpOnly: true,
			SessionFolder:   "/",
			SessionSecret:   "secret",
			CookiesSecret:   "cookie",
		})
		authService := services.NewAuthService(userService, sessionManager)

		authController := controllers.NewAuthController(authService)
		http.HandleFunc("/register", authController.Register)
		server := httptest.NewServer(http.DefaultServeMux)
		defer server.Close()

		registerDto := dtos.RegisterDto{
			Name:           "Test User",
			Email:          "test@example.com",
			Password:       "password123",
			PasswordRepeat: "password123",
		}
		payload, err := json.Marshal(registerDto)
		assert.NoError(t, err, "Failed to marshal the register DTO")

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		authController.Register(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code, "Expected HTTP status 200 OK")
		assert.Contains(t, rec.Body.String(), "User registered successfully", "Response body does not contain expected success message")
	})

}
