package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Mixturka/vm-hub/internal/app/application/dtos"
	"github.com/Mixturka/vm-hub/internal/app/infrustructure/config"
)

type BaseOAuthService struct {
	BaseURL string
	options *config.BaseOAuthProviderOptions
}

func NewBaseOAuthService(baseURL string, options *config.BaseOAuthProviderOptions) BaseOAuthService {
	return BaseOAuthService{
		BaseURL: baseURL,
		options: options,
	}
}

func (bos BaseOAuthService) RedirectURL() string {
	return bos.BaseURL + "/auth/oauth/callback/" + bos.options.Name
}

func (bos BaseOAuthService) ExtractUserInfo(data map[string]interface{}) (dtos.OAuthUserDto, error) {
	dto := dtos.OAuthUserDto{}
	var ok bool

	if dto.ID, ok = data["id"].(string); !ok {
		return dto, fmt.Errorf("missing or invalid ID field")
	}
	if dto.Picture, ok = data["picture"].(string); !ok {
		return dto, fmt.Errorf("missing or invalid Picture field")
	}
	if dto.Name, ok = data["name"].(string); !ok {
		return dto, fmt.Errorf("missing or invalid Name field")
	}
	if dto.Email, ok = data["email"].(string); !ok {
		return dto, fmt.Errorf("missing or invalid Email field")
	}
	if dto.AccessToken, ok = data["access_token"].(string); !ok {
		return dto, fmt.Errorf("missing or invalid AccessToken field")
	}
	if dto.RefreshToken, ok = data["refresh_token"].(string); !ok {
		return dto, fmt.Errorf("missing or invalid RefreshToken field")
	}
	if expiresAt, ok := data["expires_at"].(float64); ok {
		dto.ExpiresAt = int64(expiresAt)
	} else {
		return dto, fmt.Errorf("missing or invalid ExpiresAt field")
	}
	dto.Provider = bos.options.Name

	return dto, nil
}

func (bos BaseOAuthService) AuthURL() string {
	query := url.Values{}
	query.Add("response_type", "code")
	query.Add("client_id", bos.options.ClientID)
	query.Add("redirect_uri", bos.RedirectURL())
	query.Add("scope", strings.Join(bos.options.Scopes, " "))
	query.Add("access_type", "offline")
	query.Add("prompt", "select_account")
	return fmt.Sprintf("%s?%s", bos.options.AuthorizeURL, query.Encode())
}

func (bos BaseOAuthService) FindUserByCode(code string) (dtos.OAuthUserDto, error) {
	tokenQuery := url.Values{}
	tokenQuery.Set("client_id", bos.options.ClientID)
	tokenQuery.Set("client_secret", bos.options.ClientSecret)
	tokenQuery.Set("redirect_uri", bos.RedirectURL())
	tokenQuery.Set("grant_type", "authorization_code")
	tokenQuery.Set("code", code)

	resp, err := http.Post(bos.options.AccessURL, "application/x-www-form-urlencoded", strings.NewReader(tokenQuery.Encode()))
	if err != nil {
		return dtos.OAuthUserDto{}, fmt.Errorf("failed to request token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dtos.OAuthUserDto{}, fmt.Errorf("failed to request token: %s", resp.Status)
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		ExpiresAt    int64  `json:"expires_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return dtos.OAuthUserDto{}, fmt.Errorf("failed to decode token response: %v", err)
	}

	userRequest, err := http.NewRequest("GET", bos.options.ProfileURL, nil)
	if err != nil {
		return dtos.OAuthUserDto{}, fmt.Errorf("failed to create user info request: %v", err)
	}
	userRequest.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

	client := &http.Client{}
	userResp, err := client.Do(userRequest)
	if err != nil {
		return dtos.OAuthUserDto{}, fmt.Errorf("failed to fetch user info: %v", err)
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		return dtos.OAuthUserDto{}, fmt.Errorf("unauthorized: could not fetch user info from %s, check the access token", bos.options.ProfileURL)
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(userResp.Request.Body).Decode(&userInfo); err != nil {
		return dtos.OAuthUserDto{}, fmt.Errorf("failed to decode user from response: %v", err)
	}

	userData, err := bos.ExtractUserInfo(userInfo)
	if err != nil {
		return dtos.OAuthUserDto{}, fmt.Errorf("failed to extract userData from decoded json: %v", err)
	}

	return dtos.OAuthUserDto{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresAt:    tokenResponse.ExpiresAt,
		Provider:     bos.options.Name,
		ID:           userData.ID,
		Picture:      userData.Picture,
		Name:         userData.Name,
		Email:        userData.Email,
	}, nil
}
