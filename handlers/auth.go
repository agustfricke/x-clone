package handlers

import (
  "fmt"
  "strings"
  "time"

  "github.com/gofiber/fiber/v2"
  "github.com/golang-jwt/jwt"
  "github.com/pquerna/otp/totp"
  "golang.org/x/crypto/bcrypt"
)


func SignUp(c *fiber.Ctx) error {
  var payload *models.SignUpInput
  db := database.DB

  if err := c.BodyParser(&payload); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
  }

  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)

  if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
  }

  newUser := models.User{
    Name:     payload.Name,
    Email:    strings.ToLower(payload.Email),
    Password: string(hashedPassword),
  }

  result := db.Create(&newUser)

  if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
    return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "User with that email already exists"})
  } else if result.Error != nil {
    return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
  }

  tokenByte := jwt.New(jwt.SigningMethodHS256)

  now := time.Now().UTC()
  claims := tokenByte.Claims.(jwt.MapClaims)
  expDuration := time.Hour * 24

  fmt.Print(newUser.ID)

  claims["sub"] = newUser.ID
  claims["exp"] = now.Add(expDuration).Unix()
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()

  tokenString, err := tokenByte.SignedString([]byte(config.Config("SECRET_KEY")))

  if err != nil {
    return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
  }

  SendEmail(tokenString, payload.Email)

  return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": &newUser}})
}


func SignIn(c *fiber.Ctx) error {
  var payload *models.SignInInput
  db := database.DB

  if err := c.BodyParser(&payload); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
  }

  var user models.User
  result := db.First(&user, "email = ?", strings.ToLower(payload.Email))
  if result.Error != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
  }

  if !user.Verified {
    return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "No estas verificado bro!"})
  }

  err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
  if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
  }
  
  if user.Otp_enabled == true {
    valid := totp.Validate(payload.Token, user.Otp_secret)
    if !valid {
      return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "status":  "fail",
        "message": "Token 2FA not valid",
      })
    }

  }

  tokenByte := jwt.New(jwt.SigningMethodHS256)

  now := time.Now().UTC()
  claims := tokenByte.Claims.(jwt.MapClaims)
  expDuration := time.Hour * 24

  claims["sub"] = user.ID
  claims["exp"] = now.Add(expDuration).Unix()
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()

  tokenString, err := tokenByte.SignedString([]byte(config.Config("SECRET_KEY")))

  if err != nil {
    return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
  }

  c.Cookie(&fiber.Cookie{
    Name:     "token",
    Value:    tokenString,
    Path:     "/",
    MaxAge:   24 * 60 * 60,
    Secure:   false,
    HTTPOnly: true,
    Domain:   "localhost",
  })

  return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "token": tokenString})
}

func Logout(c *fiber.Ctx) error {
  expired := time.Now().Add(-time.Hour * 24)
  c.Cookie(&fiber.Cookie{
    Name:    "token",
    Value:   "",
    Expires: expired,
  })
  return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

