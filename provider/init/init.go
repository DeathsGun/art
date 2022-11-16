package init

import (
	"github.com/deathsgun/art/di"
	ihkProvider "github.com/deathsgun/art/ihk/provider"
	"github.com/deathsgun/art/provider"
	redmineProvider "github.com/deathsgun/art/redmine/provider"
	"github.com/deathsgun/art/text"
	untisProvider "github.com/deathsgun/art/untis/provider"
	"github.com/rs/zerolog/log"
)

func InitializeProvider() {
	providerService := di.Instance[provider.IProviderService]("providerService")

	err := providerService.RegisterProvider(untisProvider.NewProvider())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register untis provider")
	}
	err = providerService.RegisterProvider(redmineProvider.NewProvider())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register redmine provider")
	}
	//err = providerService.RegisterProvider(jiraProvider.NewProvider())
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Failed to register jira provider")
	//}
	err = providerService.RegisterProvider(text.NewProvider())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register text provider")
	}
	err = providerService.RegisterProvider(ihkProvider.NewProvider())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register ihk provider")
	}
}
