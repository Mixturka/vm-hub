package services

import (
	"context"
	"fmt"

	"github.com/Mixturka/vm-hub/internal/app/application/interfaces"
	"github.com/Mixturka/vm-hub/internal/app/domain/entities"
	"github.com/Mixturka/vm-hub/pkg/security"
	"github.com/google/uuid"
)

type UserService struct {
	repository interfaces.UserRepository
}

func NewUserService(repository interfaces.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (us *UserService) FindByID(ctx context.Context, id string) (*entities.User, error) {
	return us.repository.GetByID(ctx, id)
}

func (us *UserService) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	return us.repository.GetByEmail(ctx, email)
}

func (us *UserService) CreateUser(ctx context.Context, email string, password string, name string,
	profilePic string, method entities.AuthMethod, isEmailVerified bool) (*entities.User, error) {
	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}
	user := &entities.User{
		ID:              uuid.NewString(),
		ProfilePicture:  profilePic,
		Name:            name,
		Email:           email,
		Password:        hashedPassword,
		Accounts:        []entities.Account{},
		IsEmailVerified: isEmailVerified,
	}
	return user, us.repository.Save(ctx, user)
}
