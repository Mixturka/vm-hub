package services

import "github.com/Mixturka/vm-hub/internal/app/infrustructure/auth"

type ProviderService struct {
	options *auth.OAuthServiceOptions
}

func NewProviderService(options *auth.OAuthServiceOptions) *ProviderService {
	for i := range options.Services {
		options.Services[i].BaseURL = options.BaseURL
	}
	return &ProviderService{
		options: options,
	}
}

func (ps *ProviderService) GetServiceByName(name string) *auth.BaseOAuthService {
	for _, service := range ps.options.Services {
		if service.Options().Name == name {
			return &service
		}
	}
	return nil
}
