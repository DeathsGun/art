package provider

import (
	"context"
	"github.com/deathsgun/art/config/model"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/report"
	"time"
)

type impl struct {
}

func (i *impl) ValidateConfig(_ context.Context, _ *model.ProviderConfig) error {
	return nil
}

func (i *impl) Id() string {
	return "PROVIDER_JIRA"
}

func (i *impl) Logo() string {
	return "logo-gradient-blue-jira.svg"
}

func (i *impl) Capabilities() []provider.Capability {
	return []provider.Capability{
		provider.Configurable,
		provider.ConfigServer,
		provider.ConfigUsername,
		provider.ConfigPassword,
		provider.TypeImport,
	}
}

func (i *impl) Import(_ context.Context, _ time.Time) ([]report.Entry, error) {
	//TODO implement me
	panic("implement me")
}

func NewProvider() provider.ImportProvider {
	return &impl{}
}
