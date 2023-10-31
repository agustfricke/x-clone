package main

import (
	"github.com/agustfricke/x-clone/database"
	"github.com/agustfricke/x-clone/handlers"
	"github.com/agustfricke/x-clone/middleware"
	"github.com/agustfricke/x-clone/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
  database.ConnectDB()

  engine := html.New("./templates", ".html")

  app := fiber.New(fiber.Config{
    Views: engine, 
    ViewsLayout: "layouts/main", 
  })

  app.Static("/", "./public")

  app.Get("/sign_up_view", views.SignUpView)
  app.Get("/sign_in_view", views.SignInView)
  app.Get("/home", middleware.DeserializeUser, views.HomeView)
  app.Get("/", views.RootView)
  app.Get("/notifications", views.NotificationView)
  app.Post("/signup", handlers.SignUp)
  app.Get("/verify/:token", handlers.VerifyEmail)
  app.Get("/get/user/", handlers.AllUsers)

  app.Post("/signin/oh", handlers.SignIn)
  app.Post("/otp/:email/:password", handlers.VerifyOTP)

	app.Get("/auth/google", handlers.AuthGoogle)
	app.Get("/auth/google/callback", handlers.CallbackGoogle) 

  app.Get("/auth/github", handlers.AuthGitHub)
  app.Get("/auth/github/callback", handlers.CallbackGitHub)

  app.Listen(":8080")
}
