package entities

import "time"

type User struct {
	ID             string
	ProfilePicture string
	Name           string
	Email          string
	Password       string
	Accounts       []Account

	IsEmailVerified    bool
	IsTwoFactorEnabled bool
	Method             AuthMethod

	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuthMethod int

const (
	Credentials AuthMethod = iota
	Google
	Yandex
)
