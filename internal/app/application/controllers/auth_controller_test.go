package controllers_test

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	"github.com/Mixturka/vm-hub/internal/app/application/controllers"
// 	"github.com/Mixturka/vm-hub/internal/app/application/dtos"
// 	"github.com/Mixturka/vm-hub/internal/app/application/services"
// 	"github.com/Mixturka/vm-hub/internal/app/config"
// 	"github.com/Mixturka/vm-hub/internal/app/infrustructure/database/postgres"
// 	"github.com/Mixturka/vm-hub/internal/app/infrustructure/session"
// 	"github.com/Mixturka/vm-hub/internal/pkg/test"
// 	"github.com/Mixturka/vm-hub/pkg/putils"
// 	"github.com/stretchr/testify/assert"
// )

// var (
// 	migrationsPath         string
// 	absoluteMigrationsPath string
// )

// func TestMain(t *testing.M) {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		log.Fatalf("Error getting current working directory: %v", err)
// 	}

// 	projRoot, err := putils.GetProjectRoot(cwd)
// 	if err != nil {
// 		log.Fatalf("Error finding project root: %v", err)
// 	}

// 	migrationsPath = test.MustGetEnv("POSTGRES_MIGRATIONS_PATH")
// 	absoluteMigrationsPath = test.GetAbsolutePath(projRoot, migrationsPath)

// 	exitCode := t.Run()
// 	os.Exit(exitCode)
// }

// // Integrational tests
// func TestRegister_Success(t *testing.T) {
// 	t.Parallel()

// 	t.Run("Test register success", func(t *testing.T) {
// 		t.Parallel()

// 		ptUtil := test.NewPostgresTestUtilWithIsolatedSchema(t)
// 		test.ApplyMigrations(ptUtil.DB().Config().ConnString(), absoluteMigrationsPath)
// 		repo := postgres.NewPostgresUserRepository(ptUtil.DB())
// 		userService := services.NewUserService(repo)

// 		util := test.NewRedisTestUtil(t)
// 		client := util.Client()
// 		rs := session.NewRedisStore(client)

// 		sessionManager := session.NewSessionManager(rs, &config.SessionOptions{})
// 		authService := services.NewAuthService(userService, sessionManager)

// 		authController := controllers.NewAuthController(authService)
// 		http.HandleFunc("/register", authController.Register)
// 		server := httptest.NewServer(nil)
// 		defer server.Close()

// 		registerDto := dtos.RegisterDto{
// 			Email:    "test@example.com",
// 			Password: "password123",
// 			Name:     "Test User",
// 		}

// 		body, _ := json.Marshal(registerDto)
// 		req := httptest.NewRequest(http.MethodPost, server.URL+"/register", bytes.NewReader(body))
// 		rec := httptest.NewRecorder()

// 		// Send the request to the server
// 		server.Config.Handler.ServeHTTP(rec, req)

// 		// Check the response
// 		assert.Equal(t, http.StatusOK, rec.Code)

// 		user, err := userService.FindByEmail(context.Background(), registerDto.Email)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, user)
// 		assert.Equal(t, registerDto.Email, user.Email+"1'")
// 		assert.Equal(t, registerDto.Name, user.Name)

// 		// Optionally: Check if the user is assigned a session by checking the session manager
// 		// This can depend on how you implement session checking (e.g., looking at session cookies or session data)
// 		// Example assertion for session:
// 		sessionData := rec.Header().Get("Set-Cookie")
// 		assert.Contains(t, sessionData, "session_id")
// 	})

// }
