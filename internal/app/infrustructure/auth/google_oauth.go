package auth

import (
	"github.com/Mixturka/vm-hub/internal/app/application/dtos"
	"github.com/Mixturka/vm-hub/internal/app/infrustructure/config"
)

type GoogleProfile struct {
	Aud             string `json:"aud"`
	Azp             string `json:"azp"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"email_verified"`
	Exp             int64  `json:"exp"`
	FamilyName      string `json:"family_name,omitempty"`
	GivenName       string `json:"given_name"`
	Hd              string `json:"hd,omitempty"`
	Iat             string `json:"iat"`
	Iss             string `json:"iss"`
	Jti             string `json:"jti,omitempty"`
	Locale          string `json:"locale,omitempty"`
	Name            string `json:"name"`
	Nbf             int64  `json:"nbf,omitempty"`
	Picture         string `json:"picture"`
	Sub             string `json:"sub"`
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token,omitempty"`
}

type GoogleProvider struct {
	base BaseOAuthService
}

func NewGoogleProvider(options config.OAuthProviderOptions) GoogleProvider {
	baseOptions := config.BaseOAuthProviderOptions{
		Name:         "google",
		AuthorizeURL: "https://accounts.google.com/o/oauth2/v2/auth",
		AccessURL:    "https://oauth2.googleapis.com/token",
		ProfileURL:   "https://www.googleapis.com/oauth2/v3/userinfo",
		Scopes:       options.Scopes,
		ClientID:     options.CliendID,
		ClientSecret: options.ClientSecret,
	}
	return GoogleProvider{
		base: NewBaseOAuthService(&baseOptions),
	}
}

func (gp GoogleProvider) ExtractUserInfo(data *GoogleProfile) (dtos.OAuthUserDto, error) {
	return gp.base.ExtractUserInfo(map[string]interface{}{
		"email":   data.Email,
		"name":    data.Name,
		"picture": data.Picture,
	})
}
