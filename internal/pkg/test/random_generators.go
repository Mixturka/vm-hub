package test

import (
	"fmt"
	"time"

	"github.com/Mixturka/vm-hub/internal/app/domain/entities"
	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func NewRandomUser() *entities.User {
	rand.Seed(uint64(time.Now().UnixNano()))

	randomAuthMethods := []entities.AuthMethod{entities.Credentials, entities.Google, entities.Yandex}

	return &entities.User{
		ID:                 uuid.NewString(),
		Name:               fmt.Sprintf("User-%s", generateRandomString(8)),
		Email:              fmt.Sprintf("%s@example.com", uuid.NewString()),
		Password:           generateRandomString(12),
		ProfilePicture:     fmt.Sprintf("https://example.com/profile/%s.jpg", generateRandomString(8)),
		Accounts:           []entities.Account{},
		IsEmailVerified:    rand.Intn(2) == 1,
		IsTwoFactorEnabled: rand.Intn(2) == 1,
		Method:             randomAuthMethods[rand.Intn(len(randomAuthMethods))],
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
	}
}
