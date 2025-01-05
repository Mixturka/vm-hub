package entities

import "time"

type Token struct {
	ID        string
	UserEmail string
	Token     string
	Type      TokenType
	ExpiresIn time.Time
}

type TokenType int

const (
	Verification TokenType = iota
	TwoFactor
	PasswordReset
)
