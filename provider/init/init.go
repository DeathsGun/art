package init

import (
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/jira"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/redmine"
	untisProvider "github.com/deathsgun/art/untis/provider"
	"github.com/rs/zerolog/log"
)

func InitializeProvider() {
	providerService := di.Instance[provider.IProviderService]("providerService")

	err := providerService.RegisterProvider(untisProvider.NewProvider())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register untis provider")
	}
	err = providerService.RegisterProvider(redmine.NewProvider())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register redmine provider")
	}
	err = providerService.RegisterProvider(jira.NewProvider())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register redmine provider")
	}
}
