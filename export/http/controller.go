package http

import (
	"errors"
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
	app.Get("/export/start-date/:provider", auth.New, HandleStartDate)
	app.All("/", func(c *fiber.Ctx) error {
		return c.Redirect("/export", fiber.StatusTemporaryRedirect)
	})
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
		if errors.Is(err, export.ErrNoImportProviderEntries) {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return err
	}
	if len(bytes) == 0 {
		return c.SendStatus(fiber.StatusNoContent)
	}
	contentType, err := exportService.GetContentType(c.UserContext(), requestDto.Provider)
	if err != nil {
		return err
	}
	c.Set("Content-Type", contentType)
	c.Set("File-Name", fmt.Sprintf("%s", utils.LeapToPreviousMonday(requestDto.Date).Format("2006-01-02")))
	return c.Status(fiber.StatusOK).Send(bytes)
}

func HandleStartDate(c *fiber.Ctx) error {
	if c.Params("provider") == "" {
		return fiber.ErrBadRequest
	}
	providerService := di.Instance[provider.IProviderService]("providerService")
	prov, ok := providerService.GetProvider(c.Params("provider"))
	if !ok {
		return fiber.ErrNotFound
	}
	switch v := prov.(type) {
	case provider.ImportProvider:
		return fiber.ErrBadRequest
	case provider.ExportProvider:
		date, err := v.GetStartDate(c.UserContext())
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(&dto.DateResponse{Date: date})
	}
	return errors.New("unknown provider type") // Should be impossible
}
