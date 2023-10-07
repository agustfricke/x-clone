package views

import "github.com/gofiber/fiber/v2"

func SignUpView(c *fiber.Ctx) error {
	return c.Render("sign_up", fiber.Map{})
}
