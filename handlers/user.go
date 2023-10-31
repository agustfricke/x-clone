package handlers

import (
	"github.com/agustfricke/x-clone/database"
	"github.com/agustfricke/x-clone/models"
	"github.com/gofiber/fiber/v2"
)

func AllUsers(c *fiber.Ctx) error {
    db := database.DB
    var user models.User

    db.Find(&user)
    return c.JSON(user)
}

func UserProfile(c *fiber.Ctx) error {
	  id := c.Params("id")
    db := database.DB
    var user models.User

    if err := db.First(&user, "ID = ?", id).Error; err != nil {
      return c.SendStatus(fiber.StatusBadRequest)
    }

    if err := db.Preload("Posts").Find(&user).Error; err != nil {
      return c.SendStatus(fiber.StatusInternalServerError)
    }

	  return c.Render("user_profile", fiber.Map{
        "User": user,
	})
}

// func UpdateProfile(c *fiber.Ctx) error {
// }

func MyUserProfile(c *fiber.Ctx) error {
    db := database.DB

	  user := c.Locals("user").(*models.User)

    if err := db.First(&user, "ID = ?", user.ID).Error; err != nil {
      return c.SendStatus(fiber.StatusBadRequest)
    }

    if err := db.Preload("Posts").Find(&user).Error; err != nil {
      return c.SendStatus(fiber.StatusInternalServerError)
    }

	  return c.Render("my_profile", fiber.Map{
        "User": user,
	})
  /*
  profile.html
        <li class="text-red-500">{{ .User.Nickname }}</li>
        <li class="text-red-500">{{ .User.ID }}</li>
        <li class="text-red-500">{{ .User.Sub }}</li>
        <br>


        {{ range .User.Tasks }}
        <div class="bg-gray-900 m-2">
            <li class="text-white">Task: {{ .Name }}</li>
            <li class="text-white">TaskID : {{ .ID }}</li>
            <li class="text-green-200">Nickname: {{ $.User.Nickname }}</li>

            <button class="bg-red-400 text-white" hx-delete="/delete/{{ .ID }}">Delete</button>
        </div>
        {{ end }}
  */

}



