package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Mixturka/vm-hub/internal/application/interfaces"
	"github.com/Mixturka/vm-hub/internal/domain/entities"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) interfaces.UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User

	userQuery := `SELECT id, profile_picture, name, email, password,
			      	is_email_verified, is_two_factor_enabled, method, created_at, updated_at
			  	  FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, userQuery, id).Scan(&user.ID, &user.ProfilePicture, &user.Name,
		&user.Email, &user.Password, &user.IsEmailVerified,
		&user.IsTwoFactorEnabled, &user.Method, &user.CreatedAt,
		&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with ID %s not found: %w", id, sql.ErrNoRows)
	} else if err != nil {
		return nil, fmt.Errorf("error fetching user by ID: %w", err)
	}

	accountsQuery := `SELECT id, user_id, type, provider, refresh_token, access_token, expires_at, created_at, updated_at
					  FROM accounts WHERE user_id = $1`

	rows, err := r.db.Query(ctx, accountsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts for user %s: %w", id, err)
	}
	defer rows.Close()

	accounts := []entities.Account{}
	for rows.Next() {
		var account entities.Account

		if err := rows.Scan(&account.ID, &account.UserID, &account.Type, &account.Provider,
			&account.RefreshToken, &account.AccessToken, &account.ExpiresAt,
			&account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning account for user %s: %w", id, err)
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over accounts for user %s: %w", id, err)
	}

	user.Accounts = accounts
	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User

	userQuery := `SELECT id, profile_picture, name, email, password,
			      	is_email_verified, is_two_factor_enabled, method, created_at, updated_at
			  	  FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, userQuery, email).Scan(&user.ID, &user.ProfilePicture, &user.Name,
		&user.Email, &user.Password, &user.IsEmailVerified,
		&user.IsTwoFactorEnabled, &user.Method, &user.CreatedAt,
		&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user with email %s not found: %w", email, sql.ErrNoRows)
	} else if err != nil {
		return nil, fmt.Errorf("error fetching user by email: %w", err)
	}

	accountsQuery := `SELECT id, user_id, type, provider, refresh_token, access_token, expires_at, created_at, updated_at
					  FROM accounts WHERE user_id = $1`

	rows, err := r.db.Query(ctx, accountsQuery, user.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts for user with email %s: %w", email, err)
	}
	defer rows.Close()

	accounts := []entities.Account{}
	for rows.Next() {
		var account entities.Account

		if err := rows.Scan(&account.ID, &account.UserID, &account.Type, &account.Provider,
			&account.RefreshToken, &account.AccessToken, &account.ExpiresAt,
			&account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning account for user with email %s: %w", email, err)
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over accounts for user with email %s: %w", email, err)
	}

	user.Accounts = accounts
	return &user, nil
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (id, profile_picture, name, email, password,
			    is_email_verified, is_two_factor_enabled, method, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.Query(ctx, query, user.ID, user.ProfilePicture, user.Name, user.Email,
		user.Password, user.IsEmailVerified, user.IsTwoFactorEnabled, user.Method,
		user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	for _, account := range user.Accounts {
		accountQuery := `INSERT INTO accounts (id, user_id, type, provider, refresh_token, access_token, expires_at, created_at, updated_at)
						 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
		_, err := r.db.Exec(ctx, accountQuery, account.ID, user.ID, account.Type, account.Provider,
			account.RefreshToken, account.AccessToken, account.ExpiresAt,
			account.CreatedAt, account.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `UPDATE users SET profile_picture = $2, name = $3, email = $4,
			 	password = $5, is_email_verified = $6, is_two_factor_enabled = $7,
				method = $8, created_at = $9, updated_at = $10
			  WHERE id = $1`
	_, err := r.db.Exec(ctx, query, user.ID, user.ProfilePicture, user.Name, user.Email,
		user.Password, user.IsEmailVerified, user.IsTwoFactorEnabled, user.Method,
		user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	accountsQuery := "DELETE FROM accounts WHERE user_id = $1"
	_, err := r.db.Exec(ctx, accountsQuery, id)
	if err != nil {
		return err
	}

	userQuery := "DELETE FROM users WHERE id = $1"
	_, err = r.db.Exec(ctx, userQuery, id)
	return err
}
