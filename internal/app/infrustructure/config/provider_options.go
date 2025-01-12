package config

type BaseOAuthProviderOptions struct {
	Name         string
	AuthorizeURL string
	AccessURL    string
	ProfileURL   string
	Scopes       []string
	ClientID     string
	ClientSecret string
}

type OAuthProviderOptions struct {
	Scopes       []string
	CliendID     string
	ClientSecret string
}
