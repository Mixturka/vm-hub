package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Mixturka/vm-hub/internal/app/application/interfaces"
	"github.com/Mixturka/vm-hub/internal/pkg/test"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
)

var (
	testUtil *test.PostgresTestUtil
	repo     interfaces.UserRepository
)

func TestMain(m *testing.M) {
	testUtil = test.NewPostgresTestUtil()
	repo = NewPostgresUserRepository(testUtil.DB)
	exitCode := m.Run()
	defer testUtil.Close()
	os.Exit(exitCode)
}

func startTestTransaction(t *testing.T) (*test.PostgresTestUtil, context.Context, context.CancelFunc) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := testUtil.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	return testUtil, ctx, cancel
}

func prettyLog(t *testing.T, testName string, message string) {
	t.Logf("==== %s ====\n%s\n", testName, message)
}

func TestPostgresUserRepository_Save_GetByEmail(t *testing.T) {
	t.Run("Save User And Get By Email Test", func(t *testing.T) {
		t.Parallel()
		prettyLog(t, "TestPostgresUserRepository_Save_GetByEmail", "Starting test to get and save a user")

		user := *test.NewRandomUser()

		prettyLog(t, "TestPostgresUserRepository_Save_GetByEmail", "Inserting user to database")

		_, ctx, cancel := startTestTransaction(t)
		defer cancel()
		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		fetchedUser, err := repo.GetByEmail(ctx, user.Email)
		assert.NoError(t, err, "GetByEmail shouldn't return an error")
		assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
		test.AssertUsersEqual(t, &user, fetchedUser)

		prettyLog(t, "TestPostgresUserRepository_Save_GetByEmail", "User saved and verified successfully")
	})
}

func TestPostgresUserRepository_Save_GetByID(t *testing.T) {
	t.Run("Save User And Get By ID Test", func(t *testing.T) {
		t.Parallel()
		prettyLog(t, "TestPostgresUserRepository_Save_GetByID", "Starting test to save and get a user")

		user := *test.NewRandomUser()

		prettyLog(t, "TestPostgresUserRepository_Save_GetByID", "Inserting user to database")

		_, ctx, cancel := startTestTransaction(t)
		defer cancel()

		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		fetchedUser, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err, "GetByID shouldn't return an error")
		assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
		test.AssertUsersEqual(t, &user, fetchedUser)

		prettyLog(t, "TestPostgresUserRepository_Save_GetByID", "User saved and verified successfully")
	})
}

func TestPostgresUserRepository_Save_Update(t *testing.T) {
	t.Run("Save And Update User Test", func(t *testing.T) {
		t.Parallel()
		prettyLog(t, "TestPostgresUserRepository_Save_Update", "Starting test to save and update a user")

		user := *test.NewRandomUser()

		prettyLog(t, "TestPostgresUserRepository_Save_Update", "Inserting user to database")

		_, ctx, cancel := startTestTransaction(t)
		defer cancel()

		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		user.Name = "Andrew"
		err = repo.Update(ctx, &user)
		assert.NoError(t, err, "Update shouldn't return an error")

		fetchedUser, err := repo.GetByID(ctx, user.ID)
		assert.NoError(t, err, "GetByID shouldn't return an error")
		assert.NotNil(t, fetchedUser, "Fetched user should not be nil")
		test.AssertUsersEqual(t, &user, fetchedUser)

		prettyLog(t, "TestPostgresUserRepository_Save_Update", "User saved and updated successfully")
	})
}

func TestPostgresUserRepository_Save_Delete(t *testing.T) {
	t.Run("Save And Delete User Test", func(t *testing.T) {
		t.Parallel()
		prettyLog(t, "TestPostgresUserRepository_Save_Delete", "Starting test to save and delete a user")

		user := *test.NewRandomUser()

		prettyLog(t, "TestPostgresUserRepository_Save_Delete", "Inserting user to database")

		_, ctx, cancel := startTestTransaction(t)
		defer cancel()

		err := repo.Save(ctx, &user)
		assert.NoError(t, err, "Save shouldn't return an error")

		err = repo.Delete(ctx, user.ID)
		assert.NoError(t, err, "Delete shouldn't return an error")

		time.Sleep(500 * time.Millisecond)

		fetchedUser, err := repo.GetByID(ctx, user.ID)
		assert.Error(t, err, "Expected an error when fetching deleted user")
		assert.Nil(t, fetchedUser, "Fetched user should be nil after deletion")
	})
}
