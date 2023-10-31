package views

import "github.com/gofiber/fiber/v2"

func NotificationView(c *fiber.Ctx) error {
	return c.Render("notification", fiber.Map{})
}
