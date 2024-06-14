package midlewares

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"todolistserver.com/test/database"
	"todolistserver.com/test/models"
	"todolistserver.com/test/types"
)

func PrepeareProjectDefault(c *fiber.Ctx) error {
	user := c.Locals("user").(string)

	project := models.Project{
		Title: *models.GetDefaultProjectTitle(user),
		User:  user,
	}

	var counts int64
	database.DB.Model(&models.Project{}).Where(project).Count(&counts)

	if counts <= 0 {
		if err := database.DB.Create(&project).Error; err != nil {
			log.Println("Error Project: ", err)

			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{ // ANY UNEXPECTED ERR
				"msg":      types.ERR_UNEXPECTED,
				"err_type": types.ERR_TYPE_MESSAGE,
			})
		}
	}

	return c.Next()
}
