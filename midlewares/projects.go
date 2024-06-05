package midlewares

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"todolistserver.com/test/database"
	"todolistserver.com/test/models"
	"todolistserver.com/test/types"
)

func RoutesGetProjectByIdMildware(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		log.Println("Error Project: Bad Request")

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      "Invalid project ID",
		})
	}

	project := models.Project{
		ID: uint(id),
	}

	if err := database.DB.Where(project).First(&project).Error; err != nil {
		log.Println("Error Project: ", err)

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      fmt.Sprintf(types.ERR_MSG_NOT_FOUND, "project"),
		})
	}

	c.Locals("project", project)

	return c.Next()
}
