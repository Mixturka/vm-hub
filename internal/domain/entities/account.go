package entities

import "time"

type Account struct {
	ID       string
	Type     string
	Provider string
	User     User
	UserID   string

	RefreshToken string
	AccessToken  string
	ExpiresAt    int

	CreatedAt time.Time
	UpdatedAt time.Time
}
