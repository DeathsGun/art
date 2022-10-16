package http

import (
	"errors"
	"github.com/deathsgun/art/auth"
	"github.com/deathsgun/art/config"
	"github.com/deathsgun/art/config/http/dto"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/provider"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
)

func Initialize(app *fiber.App) {
	app.Get("/config", auth.New, HandleConfigView)
	app.Get("/config/:provider", auth.New, HandleGetConfig)
	app.Post("/config/:provider", auth.New, HandleSaveConfig)
	app.Delete("/config/:provider", auth.New, HandleDeleteConfig)
}

func HandleConfigView(c *fiber.Ctx) error {
	providerService := di.Instance[provider.IProviderService]("providerService")
	readyProviders, err := providerService.GetReadyProviders(c.UserContext())
	if err != nil {
		return err
	}
	unpreparedProviders, err := providerService.GetProviderWithMissingConfig(c.UserContext())
	if err != nil {
		return err
	}
	configurableProviders, err := providerService.GetConfigurableProviders(c.UserContext())
	if err != nil {
		return err
	}
	return c.Render("config/overview", struct {
		Lang                  string
		ReadyProviders        []provider.Provider
		UnpreparedProviders   []provider.Provider
		ConfigurableProviders []provider.Provider
		IsImportProvider      func(prov provider.Provider) bool
		HasCapability         func(prov string, capability string) bool
	}{
		Lang:                  c.Get("Accept-Language", "en"),
		ReadyProviders:        readyProviders,
		UnpreparedProviders:   unpreparedProviders,
		ConfigurableProviders: configurableProviders,
		IsImportProvider: func(prov provider.Provider) bool {
			for _, capability := range prov.Capabilities() {
				if capability == provider.TypeImport {
					return true
				}
			}
			return false
		},
	})
}

func HandleGetConfig(c *fiber.Ctx) error {
	prov := c.Params("provider")
	if prov == "" {
		return fiber.ErrBadRequest
	}
	configService := di.Instance[config.IConfigService]("configService")
	conf, err := configService.GetConfig(c.UserContext(), prov)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.ErrNotFound
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(dto.ToDto(conf))
}

func HandleSaveConfig(c *fiber.Ctx) error {
	prov := c.Params("provider")
	if prov == "" {
		return fiber.ErrBadRequest
	}

	conf := &dto.ProviderConfig{}
	err := c.BodyParser(conf)
	if err != nil {
		return err
	}
	configService := di.Instance[config.IConfigService]("configService")
	configModel := dto.ToModel(conf)
	configModel.Provider = strings.ToUpper(prov)

	providerService := di.Instance[provider.IProviderService]("providerService")
	err = providerService.ValidateConfig(c.UserContext(), configModel)

	if err != nil {
		return err
	}

	if err = configService.SaveProviderConfig(c.UserContext(), configModel); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func HandleDeleteConfig(c *fiber.Ctx) error {
	prov := c.Params("provider")
	if prov == "" {
		return fiber.ErrBadRequest
	}
	configService := di.Instance[config.IConfigService]("configService")
	err := configService.DeleteConfig(c.UserContext(), prov)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
