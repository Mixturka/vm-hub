package postgres

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

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
}

func TestPostgresUserRepository_Save_GetByEmail(t *testing.T) {
	t.Run("Save User And Get By Email Test", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping test in short mode...")
		}
		t.Parallel()

		ptUtil := test.NewPostgresTestUtilWithIsolatedSchema(t)
		test.ApplyMigrations(ptUtil.DB().Config().ConnString(), absoluteMigrationsPath)
		repo := NewPostgresUserRepository(ptUtil.DB())
		user := *test.NewRandomUser()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Save user
		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		// Fetch user by email
		fetchedUser, err := repo.GetByEmail(ctx, user.Email)
		assert.NoError(t, err, "GetByEmail shouldn't return an error")
		assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
		test.AssertUsersEqual(t, &user, fetchedUser)
	})
}

func TestPostgresUserRepository_Save_GetByID(t *testing.T) {

	t.Run("Save User And Get By ID Test", func(t *testing.T) {
		t.Parallel()
		ptUtil := test.NewPostgresTestUtilWithIsolatedSchema(t)
		test.ApplyMigrations(ptUtil.DB().Config().ConnString(), absoluteMigrationsPath)
		repo := NewPostgresUserRepository(ptUtil.DB())
		user := *test.NewRandomUser()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Save user
		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		// Fetch user by ID
		fetchedUser, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err, "GetByID shouldn't return an error")
		assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
		test.AssertUsersEqual(t, &user, fetchedUser)
	})
}

func TestPostgresUserRepository_Save_Update(t *testing.T) {
	t.Run("Save And Update User Test", func(t *testing.T) {
		t.Parallel()
		// prettyLog(t, "TestPostgresUserRepository_Save_Update", "Starting test to save and update a user")
		ptUtil := test.NewPostgresTestUtilWithIsolatedSchema(t)
		test.ApplyMigrations(ptUtil.DB().Config().ConnString(), absoluteMigrationsPath)
		repo := NewPostgresUserRepository(ptUtil.DB())
		user := *test.NewRandomUser()

		// prettyLog(t, "TestPostgresUserRepository_Save_Update", "Inserting user to database")
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
		test.AssertUsersEqual(t, &user, fetchedUser)

		// prettyLog(t, "TestPostgresUserRepository_Save_Update", "User saved and updated successfully")
	})
}

func TestPostgresUserRepository_Save_Delete(t *testing.T) {
	t.Run("Save And Delete User Test", func(t *testing.T) {
		t.Parallel()
		// prettyLog(t, "TestPostgresUserRepository_Save_Delete", "Starting test to save and delete a user")
		ptUtil := test.NewPostgresTestUtilWithIsolatedSchema(t)
		test.ApplyMigrations(ptUtil.DB().Config().ConnString(), absoluteMigrationsPath)
		repo := NewPostgresUserRepository(ptUtil.DB())
		user := *test.NewRandomUser()

		// prettyLog(t, "TestPostgresUserRepository_Save_Delete", "Inserting user to database")
		ctx1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repo.Save(ctx1, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		ctx2, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = repo.Delete(ctx2, user.ID)
		assert.NoError(t, err, "Delete shouldn't return an error")

		time.Sleep(500 * time.Millisecond)
		ctx3, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		fetchedUser, err := repo.GetByID(ctx3, user.ID)

		if fetchedUser == nil {
			t.Logf("USER NIL")
		}
		assert.Error(t, err, "Expected an error when fetching deleted user")
		assert.Nil(t, fetchedUser, "Fetched user should be nil after deletion")
	})
}
