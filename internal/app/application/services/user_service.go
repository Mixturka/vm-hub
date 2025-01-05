package services

import (
	"context"

	"github.com/Mixturka/vm-hub/internal/app/domain/entities"
)

type UserService interface {
	FindByID(ctx context.Context, id string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
	// UpdateUser(ctx context.Context, id string, updates *UpdateUserRequest) error
	DeleteUser(ctx context.Context, id string) error
}
