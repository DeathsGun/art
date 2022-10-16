package provider

import (
	"context"
	"github.com/deathsgun/art/config"
	"github.com/deathsgun/art/config/model"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/redmine"
	"github.com/deathsgun/art/report"
	"github.com/deathsgun/art/utils"
	"time"
)

type impl struct {
}

func (i *impl) ValidateConfig(ctx context.Context, conf *model.ProviderConfig) error {
	api, err := redmine.NewRedmineAPI(ctx, conf.Server, &redmine.Authorization{
		User:     conf.Username,
		Password: conf.Password,
	})
	if err != nil {
		return err
	}
	conf.Token = api.User.ApiKey
	return di.Instance[config.IConfigService]("configService").SaveProviderConfig(ctx, conf)
}

func (i *impl) Id() string {
	return "PROVIDER_REDMINE"
}

func (i *impl) Logo() string {
	return "Redmine_logo.svg"
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

func (i *impl) Import(ctx context.Context, monday time.Time) ([]report.Entry, error) {
	service := di.Instance[config.IConfigService]("configService")
	conf, err := service.GetConfig(ctx, i.Id())
	if err != nil {
		return nil, err
	}
	api, err := redmine.NewRedmineAPI(ctx, conf.Server, &redmine.Authorization{
		User:     conf.Username,
		Password: conf.Password,
	})
	endDate := monday.Add(time.Hour * 24 * 4)
	entries, err := api.GetTimeEntries(ctx, 100, 0, monday, endDate)
	if err != nil {
		return nil, err
	}
	if len(entries.TimeEntries) == 0 {
		return nil, nil
	}
	sortedEntries := map[string][]redmine.TimeEntry{}
	for _, entry := range entries.TimeEntries {
		value, ok := sortedEntries[entry.Project.Name]
		if !ok {
			value = []redmine.TimeEntry{}
		}
		value = append(value, entry)
		sortedEntries[entry.Project.Name] = value
	}
	var addedKeys []string
	var result []report.Entry
	for key, sameTimes := range sortedEntries {
		for _, entry := range sameTimes {
			startDateTime, err := time.Parse("2006-01-02", entry.SpentOn)
			if err != nil {
				return nil, err
			}
			if !utils.Contains(addedKeys, key) {
				result = append(result, report.Entry{
					Date:     startDateTime,
					Text:     key + ":",
					Category: report.Activity,
				})
				addedKeys = append(addedKeys, key)
			}
			result = append(result, report.Entry{
				Date:     startDateTime.Add(time.Hour * 10),
				Text:     "\t" + entry.Comments,
				Category: report.Activity,
			})
		}
	}
	return nil, nil
}

func NewProvider() provider.ImportProvider {
	return &impl{}
}
