package dtos

type OAuthUserDto struct {
	ID           string
	Picture      string
	Name         string
	Email        string
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
	Provider     string
}
