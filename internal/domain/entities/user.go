package entities

import "time"

type User struct {
	ID             string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProfilePicture string `gorm:"default:''"`
	Name           string
	Email          string `gorm:"unique"`
	Password       string
	Accounts       []Account `gorm:"foreignKey:UserID"`

	IsEmailVerified    bool `gorm:"default:false"`
	IsTwoFactorEnabled bool `gorm:"default:false"`
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
