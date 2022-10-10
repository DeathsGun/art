package provider

import (
	"context"
	"errors"
	"github.com/deathsgun/art/config"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/utils"
	"strings"
)

var (
	ErrProviderAlreadyRegistered = errors.New("provider already registered")
)

type IProviderService interface {
	RegisterProvider(provider Provider) error
	GetReadyProviders(ctx context.Context) ([]Provider, error)
	GetProviderWithMissingConfig(ctx context.Context) ([]Provider, error)
	GetProvider(params string) (Provider, bool)
	GetConfigurableProviders(ctx context.Context) ([]Provider, error)
}

type service struct {
	provider []Provider
}

func (s *service) GetProvider(id string) (Provider, bool) {
	for _, provider := range s.provider {
		if strings.ToLower(id) == strings.ToLower(provider.Id()) {
			return provider, true
		}
	}
	return nil, false
}

func (s *service) RegisterProvider(provider Provider) error {
	for _, p := range s.provider {
		if p.Id() == provider.Id() {
			return ErrProviderAlreadyRegistered
		}
	}
	s.provider = append(s.provider, provider)
	return nil
}

func (s *service) GetReadyProviders(ctx context.Context) ([]Provider, error) {
	configService := di.Instance[config.IConfigService]("configService")
	providerNames, err := configService.GetConfiguredProviderNames(ctx)
	if err != nil {
		return nil, err
	}
	var result []Provider
	for _, p := range s.provider {
		if utils.Contains(providerNames, p.Id()) || !utils.Contains(p.Capabilities(), Configurable) {
			result = append(result, p)
		}
	}
	return result, nil
}

func (s *service) GetConfigurableProviders(_ context.Context) ([]Provider, error) {
	var result []Provider
	for _, p := range s.provider {
		if utils.Contains(p.Capabilities(), Configurable) {
			result = append(result, p)
		}
	}
	return result, nil
}

func (s *service) GetProviderWithMissingConfig(ctx context.Context) ([]Provider, error) {
	configService := di.Instance[config.IConfigService]("configService")
	providerNames, err := configService.GetConfiguredProviderNames(ctx)
	if err != nil {
		return nil, err
	}
	var result []Provider
	for _, p := range s.provider {
		if !utils.Contains(providerNames, p.Id()) {
			if !utils.Contains(p.Capabilities(), Configurable) { // Only show provider with config
				continue
			}
			result = append(result, p)
		}
	}
	return result, nil
}

func New() IProviderService {
	return &service{provider: []Provider{}}
}
