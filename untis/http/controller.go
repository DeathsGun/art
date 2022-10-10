package http

import (
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/untis"
	"github.com/deathsgun/art/untis/http/dto"
	"github.com/gofiber/fiber/v2"
)

func Initialize(app *fiber.App) {
	app.Post("/untis/search", HandleSearchSchool)
}

func HandleSearchSchool(c *fiber.Ctx) error {
	searchRequest := &dto.SearchRequest{}
	err := c.BodyParser(searchRequest)
	if err != nil {
		return err
	}

	service := di.Instance[untis.IUntisService]("untis")
	schools, err := service.SearchSchools(c.Context(), searchRequest.Search)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(schools)
}
