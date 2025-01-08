package dtos

type RegisterDto struct {
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6"`
	PasswordRepeat string `json:"password_repeat" validate:"required,eqfield=Password"`
}
