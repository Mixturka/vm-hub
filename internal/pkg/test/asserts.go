package test

import (
	"testing"
	"time"

	"github.com/Mixturka/vm-hub/internal/app/domain/entities"
	"github.com/stretchr/testify/assert"
)

func AssertUsersEqual(t *testing.T, expected, actual *entities.User) {
	// Compare the fields that don't involve time
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.ProfilePicture, actual.ProfilePicture)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Password, actual.Password)
	assert.Equal(t, expected.Accounts, actual.Accounts)
	assert.Equal(t, expected.IsEmailVerified, actual.IsEmailVerified)
	assert.Equal(t, expected.IsTwoFactorEnabled, actual.IsTwoFactorEnabled)
	assert.Equal(t, expected.Method, actual.Method)

	// Compare time fields with a tolerance
	assert.True(t, actual.CreatedAt.Sub(expected.CreatedAt) < time.Millisecond, "CreatedAt should match within a millisecond")
	assert.True(t, actual.UpdatedAt.Sub(expected.UpdatedAt) < time.Millisecond, "UpdatedAt should match within a millisecond")
}
