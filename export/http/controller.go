package http

import (
	"fmt"
	"github.com/deathsgun/art/auth"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/export"
	"github.com/deathsgun/art/export/http/dto"
	"github.com/deathsgun/art/provider"
	"github.com/deathsgun/art/utils"
	"github.com/gofiber/fiber/v2"
)

func Initialize(app *fiber.App) {
	app.Get("/export", auth.New, HandleExportView)
	app.Post("/export", auth.New, HandleExport)
}

func HandleExportView(c *fiber.Ctx) error {
	providerService := di.Instance[provider.IProviderService]("providerService")
	providers, err := providerService.GetReadyProviders(c.UserContext())
	if err != nil {
		return err
	}
	var exportProviders []string
	for _, p := range providers {
		switch v := p.(type) {
		case provider.ExportProvider:
			exportProviders = append(exportProviders, v.Id())
		}
	}
	return c.Render("export/index", struct {
		Lang      string
		Providers []string
	}{
		Lang:      c.Get("Accept-Language", "en"),
		Providers: exportProviders,
	})
}

func HandleExport(c *fiber.Ctx) error {
	requestDto := &dto.ExportRequest{}
	if err := c.BodyParser(requestDto); err != nil {
		return err
	}
	exportService := di.Instance[export.IExportService]("exportService")
	bytes, err := exportService.Export(c.UserContext(), requestDto.Provider, requestDto.Date)
	if err != nil {
		return err
	}
	contentType, err := exportService.GetContentType(c.UserContext(), requestDto.Provider)
	if err != nil {
		return err
	}
	c.Set("Content-Type", contentType)
	c.Set("File-Name", fmt.Sprintf("%s", utils.LeapToPreviousMonday(requestDto.Date).Format("2006-01-02")))
	return c.Status(fiber.StatusOK).Send(bytes)
}
