package postgres_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Mixturka/vm-hub/internal/app/infrustructure/database/postgres"
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
		repo := postgres.NewPostgresUserRepository(ptUtil.DB())
		user := *test.NewRandomUser()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		fetchedUser, err := repo.GetByEmail(ctx, user.Email)
		assert.NoError(t, err, "GetByEmail shouldn't return an error")
		assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
		test.AssertUsersEqual(t, &user, fetchedUser)
	})
}

func TestPostgresUserRepository_Save_GetByID(t *testing.T) {

	t.Run("Save User And Get By ID Test", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping test in short mode...")
		}
		t.Parallel()
		ptUtil := test.NewPostgresTestUtilWithIsolatedSchema(t)
		test.ApplyMigrations(ptUtil.DB().Config().ConnString(), absoluteMigrationsPath)
		repo := postgres.NewPostgresUserRepository(ptUtil.DB())
		user := *test.NewRandomUser()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		fetchedUser, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err, "GetByID shouldn't return an error")
		assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
		test.AssertUsersEqual(t, &user, fetchedUser)
	})
}

func TestPostgresUserRepository_Save_Update(t *testing.T) {
	t.Run("Save And Update User Test", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping test in short mode...")
		}
		t.Parallel()
		ptUtil := test.NewPostgresTestUtilWithIsolatedSchema(t)
		test.ApplyMigrations(ptUtil.DB().Config().ConnString(), absoluteMigrationsPath)
		repo := postgres.NewPostgresUserRepository(ptUtil.DB())
		user := *test.NewRandomUser()

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
	})
}

func TestPostgresUserRepository_Save_Delete(t *testing.T) {
	t.Run("Save And Delete User Test", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping test in short mode...")
		}
		t.Parallel()

		ptUtil := test.NewPostgresTestUtilWithIsolatedSchema(t)
		test.ApplyMigrations(ptUtil.DB().Config().ConnString(), absoluteMigrationsPath)
		repo := postgres.NewPostgresUserRepository(ptUtil.DB())
		user := *test.NewRandomUser()

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
