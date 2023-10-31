package views

import "github.com/gofiber/fiber/v2"

func SignUpView(c *fiber.Ctx) error {
	return c.Render("sign_up", fiber.Map{})
}

func RootView(c *fiber.Ctx) error {
	return c.Render("root", fiber.Map{})
}

func SignInView(c *fiber.Ctx) error {
	return c.Render("sign_in", fiber.Map{})
}

func HomeView(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{})
}
