package postgres

import (
	"context"
	"database/sql"
	"vm-hub/internal/application/interfaces"
	"vm-hub/internal/domain/entities"

	"github.com/jackc/pgx/v4"
)

type PostgresUserRepository struct {
	db *pgx.Conn
}

func NewPostgresUserRepository(db *pgx.Conn) interfaces.UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, profile_picture, name, email, password,
			    is_email_verified, is_two_factor_enabled, method, created_at, updated_at
			  FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(user.ID, user.ProfilePicture, user.Name,
		user.Email, user.Password, user.IsEmailVerified,
		user.IsTwoFactorEnabled, user.Method, user.CreatedAt,
		user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, profile_picture, name, email, password,
			    is_email_verified, is_two_factor_enabled, method, created_at, updated_at
			  FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(user.ID, user.ProfilePicture, user.Name,
		user.Email, user.Password, user.IsEmailVerified,
		user.IsTwoFactorEnabled, user.Method, user.CreatedAt,
		user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (id, profile_picture, name, email, password,
			    is_email_verified, is_two_factor_enabled, method, created_at, updated_at)
			  VALUES $1, $2, $3, $4, $5, $6, $7, $8, $9, $10`
	_, err := r.db.Query(ctx, query, user.ID, user.ProfilePicture, user.Name, user.Email,
		user.Password, user.IsEmailVerified, user.IsTwoFactorEnabled, user.Method,
		user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `UPDATE users SET profile_picture = $2, name = $3, email = $4,
			 	password = $5, is_email_verified = $6, is_two_factor_enabled = $7,
				method = $8, created_at = $9, updated_at = $10
			  WHERE id = $1`
	_, err := r.db.Query(ctx, query, user.ID, user.ProfilePicture, user.Name, user.Email,
		user.Password, user.IsEmailVerified, user.IsTwoFactorEnabled, user.Method,
		user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Query(ctx, query, id)
	return err
}
