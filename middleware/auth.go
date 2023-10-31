package middleware

import (
	"fmt"
	"strings"

	"github.com/agustfricke/x-clone/config"
	"github.com/agustfricke/x-clone/database"
	"github.com/agustfricke/x-clone/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func DeserializeUser(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
    c.Redirect("/")
    return nil
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

    return []byte(config.Config("SECRET_KEY")), nil
	})

	if err != nil {
    c.Redirect("/")
    return nil
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
    c.Redirect("/")
    return nil
	}

	var user models.User
  db := database.DB
	db.First(&user, "id = ?", fmt.Sprint(claims["sub"]))

	if float64(user.ID) != claims["sub"] {
    c.Redirect("/")
    return nil
	}

	c.Locals("user", &user)

	return c.Next()
}
