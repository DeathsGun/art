package provider

import (
	"context"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/report"
	"time"
)

type impl struct {
}

func (i *impl) Logo() string {
	return "WebUntis_Logo_orange.png"
}

func (i *impl) Id() string {
	return "PROVIDER_UNTIS"
}

func (i *impl) Capabilities() []provider.Capability {
	return []provider.Capability{
		provider.TypeImport,
	}
}

func NewProvider() provider.ImportProvider {
	return &impl{}
}

func (i *impl) Import(ctx context.Context, monday time.Time) ([]report.Entry, error) {
	return nil, nil
}
