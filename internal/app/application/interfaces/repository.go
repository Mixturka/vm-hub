package interfaces

import (
	"context"

	"github.com/Mixturka/vm-hub/internal/app/domain/entities"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Save(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id string) error
}
