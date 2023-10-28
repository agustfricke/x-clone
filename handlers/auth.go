package handlers

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
	"time"

	"github.com/agustfricke/x-clone/config"
	"github.com/agustfricke/x-clone/database"
	"github.com/agustfricke/x-clone/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)


func SignUp(c *fiber.Ctx) error {

  name := c.FormValue("name")
  email     := c.FormValue("email")
  password  := c.FormValue("password")

  db := database.DB

  if name == "" || email == "" || password == "" {
    return c.SendStatus(fiber.StatusBadRequest)
  }

  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

  if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
  }

  newUser := models.User{
    Name: name, 
    Email:  strings.ToLower(email),
    Password: string(hashedPassword),
  }

  result := db.Create(&newUser)

  if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
    return c.SendStatus(fiber.StatusConflict)
  } else if result.Error != nil {
    return c.SendStatus(fiber.StatusBadGateway)
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

  SendEmail(tokenString, email)

  return c.SendStatus(fiber.StatusCreated)
}


func SendEmail(token string, email string) {
	secretPassword := config.Config("EMAIL_SECRET_KEY")
	host := config.Config("HOST")
	auth := smtp.PlainAuth(
		"",
		"agustfricke@gmail.com",
		secretPassword,
		"smtp.gmail.com",
	)

	tmpl, err := template.ParseFiles("templates/verify_email.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	data := struct {
		Token string
		Host  string
	}{
		Token: token,
		Host:  host,
	}

	var bodyContent bytes.Buffer
	err = tmpl.Execute(&bodyContent, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	content := fmt.Sprintf("To: %s\r\n"+
		"Subject: Verify Your Email Address\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"\r\n"+
		"%s", email, bodyContent.String())

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"agustfricke@gmail.com",
		[]string{email},
		[]byte(content),
	)
	if err != nil {
		fmt.Println(err)
	}
}


func VerifyEmail(c *fiber.Ctx) error {
  tokenString := c.Params("token")

  if tokenString == "" {
    return c.SendStatus(fiber.StatusUnauthorized)
  }

  tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
    if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
    }
    return []byte(config.Config("SECRET_KEY")), nil
  })

  if err != nil {
    return c.SendStatus(fiber.StatusUnauthorized)
  }

  claims, ok := tokenByte.Claims.(jwt.MapClaims)
  if !ok || !tokenByte.Valid {
    return c.SendStatus(fiber.StatusUnauthorized)
  }

	var user models.User
  db := database.DB
	db.First(&user, "id = ?", fmt.Sprint(claims["sub"]))

  if float64(user.ID) != claims["sub"] {
    return c.SendStatus(fiber.StatusForbidden)
  }

  user.Verified = true

  if err := db.Save(&user).Error; err != nil {
    return c.SendStatus(fiber.StatusInternalServerError)
  }

  return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}



