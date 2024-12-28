package entities

import "time"

type Account struct {
	ID       string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Type     string
	Provider string
	User     User `gorm:"constraint:OnDelete:CASCADE"`
	UserID   string

	RefreshToken string
	AccessToken  string
	ExpiresAt    int

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
