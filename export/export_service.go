package export

import (
	"context"
	"errors"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/report"
	"github.com/deathsgun/art/utils"
	"time"
)

var (
	ErrNoImportProviderEntries = errors.New("got no entries from import provider")
)

type IExportService interface {
	Export(ctx context.Context, prov string, date time.Time) ([]byte, error)
	GetContentType(ctx context.Context, prov string) (string, error)
}

type service struct {
}

func (s *service) GetContentType(ctx context.Context, prov string) (string, error) {
	providerService := di.Instance[provider.IProviderService]("providerService")
	providers, err := providerService.GetReadyProviders(ctx)
	if err != nil {
		return "", err
	}
	var export provider.ExportProvider
	for _, p := range providers {
		switch v := p.(type) {
		case provider.ExportProvider:
			if v.Id() == prov {
				export = v
			}
		}
	}
	if export == nil {
		return "", provider.ErrProviderNotFound
	}
	return export.ContentType(), nil
}

func (s *service) Export(ctx context.Context, prov string, date time.Time) ([]byte, error) {
	providerService := di.Instance[provider.IProviderService]("providerService")
	providers, err := providerService.GetReadyProviders(ctx)
	if err != nil {
		return nil, err
	}
	var export provider.ExportProvider
	var importProviders []provider.ImportProvider
	for _, p := range providers {
		switch v := p.(type) {
		case provider.ImportProvider:
			importProviders = append(importProviders, v)
		case provider.ExportProvider:
			if v.Id() == prov {
				export = v
			}
		}
	}
	if export == nil {
		return nil, provider.ErrProviderNotFound
	}
	if date.IsZero() {
		d, err := export.GetStartDate(ctx)
		if err != nil {
			return nil, err
		}
		date = d
	}
	monday := utils.LeapToPreviousMonday(date)
	rep := &report.Report{}
	for _, importProvider := range importProviders {
		entries, err := importProvider.Import(ctx, monday)
		if err != nil {
			return nil, err
		}
		rep.Entries = append(rep.Entries, entries...)
	}
	if len(rep.Entries) == 0 {
		return nil, ErrNoImportProviderEntries
	}
	return export.Export(ctx, rep)
}

func New() IExportService {
	return &service{}
}
