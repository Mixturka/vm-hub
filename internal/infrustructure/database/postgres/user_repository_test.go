package postgres

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Mixturka/vm-hub/internal/application/interfaces"
	"github.com/Mixturka/vm-hub/internal/domain/entities"
	"github.com/Mixturka/vm-hub/pkg/putils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db *pgxpool.Pool
var repo interfaces.UserRepository

func TestMain(m *testing.M) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	log.Printf("Current working directory: %s", cwd)

	projRoot, err := putils.GetProjectRoot(cwd)
	if err != nil {
		log.Fatalf("Error looking for root directory: %v", err)
	}

	if err := godotenv.Load(projRoot + "/.env.test"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connStr := os.Getenv("TEST_POSTGRES_URL")
	if connStr == "" {
		log.Fatal("TEST_POSTGRES_URL is not set in environment")
	}

	migrationsPath := os.Getenv("POSTGRES_MIGRATIONS_PATH")
	if migrationsPath == "" {
		log.Fatal("POSTGRES_MIGRATIONS_PATH is not set in environment")
	}

	migration, err := migrate.New("file://"+filepath.Join(projRoot, migrationsPath), connStr)
	if err != nil {
		log.Fatalf("Failed to initialize golang-migrate: %v", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	db, err = pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Error creating connection pool: %v", err)
	}

	db.Config().MaxConns = 10
	repo = NewPostgresUserRepository(db)
	exitCode := m.Run()

	log.Print("Tests run. Rolling back migrations...")
	if err := migration.Down(); err != nil {
		log.Fatalf("Failed to rollback migrations: %v", err)
	}

	defer db.Close()

	os.Exit(exitCode)
}

func truncateTables(t *testing.T) {
	_, err := db.Exec(context.Background(), `TRUNCATE users, accounts RESTART IDENTITY CASCADE;`)
	if err != nil {
		t.Fatalf("Failed to truncate tables: %v", err)
	}
}

func prettyLog(t *testing.T, testName string, message string) {
	t.Logf("==== %s ====\n%s\n", testName, message)
}

func setupDB(t *testing.T) {
	truncateTables(t)
}

func newTestUser() *entities.User {
	fixedTime := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	return &entities.User{
		ID:                 uuid.NewString(),
		Name:               "John",
		Email:              "test@email.com",
		Password:           "hashedpass",
		ProfilePicture:     "none",
		Accounts:           []entities.Account{},
		IsEmailVerified:    true,
		IsTwoFactorEnabled: true,
		Method:             entities.Yandex,
		CreatedAt:          fixedTime,
		UpdatedAt:          fixedTime,
	}
}

func TestPostgresUserRepository_Save_GetByEmail(t *testing.T) {
	t.Run("Save User And Get By Email Test", func(t *testing.T) {
		t.Parallel()
		prettyLog(t, "TestPostgresUserRepository_Save_GetByEmail", "Starting test to get and save a user")

		setupDB(t)

		user := *newTestUser()

		prettyLog(t, "TestPostgresUserRepository_Save_GetByEmail", "Inserting user to database")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		fetchedUser, err := repo.GetByEmail(context.Background(), user.Email)
		assert.NoError(t, err, "GetByEmail shouldn't return an error")
		assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
		assert.Equal(t, user, *fetchedUser, "Users should match")

		prettyLog(t, "TestPostgresUserRepository_Save_GetByEmail", "User saved and verified successfully")
	})
}

func TestPostgresUserRepository_Save_GetByID(t *testing.T) {
	t.Run("Save User And Get By ID Test", func(t *testing.T) {
		t.Parallel()
		prettyLog(t, "TestPostgresUserRepository_Save_GetByID", "Starting test to save and get a user")

		setupDB(t)

		user := *newTestUser()

		prettyLog(t, "TestPostgresUserRepository_Save_GetByID", "Inserting user to database")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		fetchedUser, err := repo.GetByID(context.Background(), user.ID)
		assert.NoError(t, err, "GetByID shouldn't return an error")
		assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
		require.Equal(t, user, *fetchedUser, "Users should match")

		prettyLog(t, "TestPostgresUserRepository_Save_GetByID", "User saved and verified successfully")
	})
}

func TestPostgresUserRepository_Save_Update(t *testing.T) {
	t.Run("Save And Update User Test", func(t *testing.T) {
		t.Parallel()
		prettyLog(t, "TestPostgresUserRepository_Save_Update", "Starting test to save and update a user")

		setupDB(t)

		user := *newTestUser()

		prettyLog(t, "TestPostgresUserRepository_Save_Update", "Inserting user to database")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		user.Name = "Andrew"
		err = repo.Update(ctx, &user)
		assert.NoError(t, err, "Update shouldn't return an error")

		fetchedUser, err := repo.GetByID(context.Background(), user.ID)
		assert.NoError(t, err, "GetByID shouldn't return an error")
		assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
		require.Equal(t, user, *fetchedUser, "Updated local user and user in db should match")

		prettyLog(t, "TestPostgresUserRepository_Save_Update", "User saved and updated successfully")
	})
}

func TestPostgresUserRepository_Save_Delete(t *testing.T) {
	t.Run("Save And Delete User Test", func(t *testing.T) {
		t.Parallel()
		prettyLog(t, "TestPostgresUserRepository_Save_Delete", "Starting test to save and delete a user")

		setupDB(t)

		user := *newTestUser()

		prettyLog(t, "TestPostgresUserRepository_Save_Update", "Inserting user to database")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		err = repo.Delete(ctx, user.ID)
		assert.NoError(t, err, "Delete shouldn't return an error")

		fetchedUser, err := repo.GetByID(ctx, user.ID)
		assert.Error(t, err, "Expected an error when fetching deleted user")
		assert.Nil(t, fetchedUser, "Fetched user should be nil after deletion")
	})
}
