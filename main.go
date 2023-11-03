package main

import (
	"log"

	"github.com/gofiber/contrib/websocket"

	"github.com/agustfricke/x-clone/database"
	"github.com/agustfricke/x-clone/handlers"
	"github.com/agustfricke/x-clone/socket"
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
  app.Get("/home", views.HomeView)
  app.Get("/", views.RootView)
  app.Get("/notifications", views.NotificationView)

  app.Get("/get/user/", handlers.AllUsers)

  app.Get("/verify/:token", handlers.VerifyEmail)

  app.Post("/signup", handlers.SignUp)
  app.Post("/signin/oh", handlers.SignIn)

  app.Post("/otp/:email/:password", handlers.VerifyOTP) // test
  app.Get("/generte/otp", handlers.GenerateOTP) // test
  app.Get("/disable/otp", handlers.DisableOTP) // test

	app.Get("/auth/google", handlers.AuthGoogle)
	app.Get("/auth/google/callback", handlers.CallbackGoogle) 
  app.Get("/auth/github", handlers.AuthGitHub)
  app.Get("/auth/github/callback", handlers.CallbackGitHub)

  app.Use(func(c *fiber.Ctx) error {
    if websocket.IsWebSocketUpgrade(c) {
      return c.Next()
    }
    return c.SendStatus(fiber.StatusUpgradeRequired)
  })

  redisClient, err := socket.InitializeRedis()
  if err != nil {
    log.Fatal("Failed to connect to Redis:", err)
  }

  socket.MessageStorage = &socket.MessageStore {
    Client: redisClient,
  }

  go socket.RunHub()

  app.Get("/ws/feed", websocket.New(socket.FeedWebsocket))

  app.Listen(":8080")
}
