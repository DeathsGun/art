package auth

import (
	"context"
	"github.com/deathsgun/art/untis"
	"github.com/gofiber/fiber/v2"
)

//TODO: Comment in
/*func New(c *fiber.Ctx) error {
	rawSession := c.Cookies("session")
	if rawSession == "" {
		if c.Accepts("application/json") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Redirect("/login?redirect="+base64.StdEncoding.EncodeToString([]byte(c.OriginalURL())), fiber.StatusTemporaryRedirect)
	}
	session, err := DecryptSession(rawSession)
	if err != nil {
		c.ClearCookie("session")
		if c.Accepts("application/json") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Redirect("/login?redirect="+base64.StdEncoding.EncodeToString([]byte(c.OriginalURL())), fiber.StatusTemporaryRedirect)
	}
	untisService := di.Instance[untis.IUntisService]("untis")
	err = untisService.ValidateLogin(c.UserContext(), session)
	if err != nil {
		c.ClearCookie("session")
		if c.Accepts("application/json") == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Redirect("/login?redirect="+base64.StdEncoding.EncodeToString([]byte(c.OriginalURL())), fiber.StatusTemporaryRedirect)
	}
	c.Locals("session", session)
	c.SetUserContext(context.WithValue(c.UserContext(), "session", session))
	return c.Next()
}*/

func New(c *fiber.Ctx) error {
	session := &untis.Session{
		Endpoint:   "https://dummy.offline",
		SessionId:  "OFFLINE",
		School:     "BK-OFFLINE",
		Token:      "OFFLINE",
		PersonType: 0,
		PersonId:   0,
	}
	c.Locals("session", session)
	c.SetUserContext(context.WithValue(c.UserContext(), "session", session))
	return c.Next()
}
