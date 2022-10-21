package http

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/deathsgun/art/auth"
	"github.com/deathsgun/art/auth/http/dto"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/untis"
	"github.com/gofiber/fiber/v2"
)

func Initialize(app *fiber.App) {
	app.Post("/login", HandleLogin)
	app.Get("/login", HandleLoginView)
}

func HandleLogin(c *fiber.Ctx) error {
	loginRequest := &dto.LoginRequest{}
	err := c.BodyParser(loginRequest)
	if err != nil {
		return err
	}
	service := di.Instance[untis.IUntisService]("untis")

	schools, err := service.SearchSchools(c.UserContext(), loginRequest.School)
	if err != nil {
		return err
	}

	var school *untis.School
	for _, s := range schools {
		if s.LoginName == loginRequest.School {
			school = &s
		}
	}
	if school == nil {
		fmt.Printf("%+v", loginRequest)
		return errors.New("school not found")
	}

	sess, err := service.Login(c.UserContext(), school, loginRequest.Username, loginRequest.Password)
	if err != nil {
		return err
	}
	session, err := auth.EncryptSession(sess)
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:  "session",
		Value: session,
		//Expires: time.Now().Add(time.Hour * 8),
	})
	data, err := base64.StdEncoding.DecodeString(c.Query("redirect", "Lw=="))
	if err != nil {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.Redirect(string(data), fiber.StatusSeeOther)
}

func HandleLoginView(c *fiber.Ctx) error {
	rawSession := c.Cookies("session")
	if rawSession != "" {
		session, err := auth.DecryptSession(rawSession)
		if err == nil {
			untisService := di.Instance[untis.IUntisService]("untis")
			err = untisService.ValidateLogin(c.UserContext(), session)
			if err == nil {
				url, err := base64.StdEncoding.DecodeString(c.Query("redirect", "Lw=="))
				if err == nil {
					return c.Redirect(string(url), fiber.StatusTemporaryRedirect)
				}
			}
		}
	}
	return c.Render("auth/login", struct {
		Lang     string
		Redirect string
	}{
		Lang:     c.Get("Accept-Language", "en"),
		Redirect: c.Query("redirect", "Lw=="),
	})
}
