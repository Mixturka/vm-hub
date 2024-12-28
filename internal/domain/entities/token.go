package entities

import "time"

type Token struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserEmail string
	Token     string `gorm:"unique"`
	Type      TokenType
	ExpiresIn time.Time
}

type TokenType int

const (
	Verification TokenType = iota
	TwoFactor
	PasswordReset
)
