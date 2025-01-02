package postgres

import (
	"context"
	"database/sql"

	"github.com/Mixturka/vm-hub/internal/application/interfaces"
	"github.com/Mixturka/vm-hub/internal/domain/entities"

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
	userQuery := `SELECT id, profile_picture, name, email, password,
			      	is_email_verified, is_two_factor_enabled, method, created_at, updated_at
			  	  FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, userQuery, id).Scan(user.ID, user.ProfilePicture, user.Name,
		user.Email, user.Password, user.IsEmailVerified,
		user.IsTwoFactorEnabled, user.Method, user.CreatedAt,
		user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	accountsQuery := `SELECT id, user_id, type, provider, refresh_token, access_token, expires_at, created_at, updated_at
					  FROM accounts WHERE user_id = $1`
	rows, err := r.db.Query(ctx, accountsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []entities.Account
	for rows.Next() {
		var account entities.Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Type, &account.Provider,
			&account.RefreshToken, &account.AccessToken, &account.ExpiresAt,
			&account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, err
		}
		account.User = user
		accounts = append(accounts, account)
	}

	user.Accounts = accounts
	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	userQuery := `SELECT id, profile_picture, name, email, password,
			    	is_email_verified, is_two_factor_enabled, method, created_at, updated_at
			      FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, userQuery, email).Scan(user.ID, user.ProfilePicture, user.Name,
		user.Email, user.Password, user.IsEmailVerified,
		user.IsTwoFactorEnabled, user.Method, user.CreatedAt,
		user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	accountsQuery := `SELECT id, user_id, type, provider, refresh_token, access_token, expires_at, created_at, updated_at
					  FROM accounts WHERE user_id = $1`
	rows, err := r.db.Query(ctx, accountsQuery, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []entities.Account
	for rows.Next() {
		var account entities.Account
		if err := rows.Scan(&account.ID, &account.UserID, &account.Type, &account.Provider,
			&account.RefreshToken, &account.AccessToken, &account.ExpiresAt,
			&account.CreatedAt, &account.UpdatedAt); err != nil {
			return nil, err
		}
		account.User = user
		accounts = append(accounts, account)
	}

	user.Accounts = accounts
	return &user, nil
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (id, profile_picture, name, email, password,
			    is_email_verified, is_two_factor_enabled, method, created_at, updated_at)
			  VALUES $1, $2, $3, $4, $5, $6, $7, $8, $9, $10`
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
