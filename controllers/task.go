package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"todolistserver.com/test/database"
	"todolistserver.com/test/models"
	"todolistserver.com/test/types"
)

func RegisterTask(c *fiber.Ctx) error {
	user := c.Locals("user").(string)
	isDefaultProject := false

	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		log.Println("Error Task: ", err)
		return c.Status(http.StatusBadRequest).SendString(types.ERR_MSG_BAR_BODY_PARSE)
	}

	project := models.Project{
		ID:   uint(task.ProjectID),
		User: user,
	}

	if project.ID == 0 {
		if err := database.DB.Where(models.Project{Title: *models.GetDefaultProjectTitle(user)}).First(&project).Error; err != nil {
			log.Println("Error Task: ", err)

			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"err_type": types.ERR_TYPE_MESSAGE,
				"msg":      types.ERR_UNEXPECTED,
			})
		}
		isDefaultProject = true
	} else if err := database.DB.Where(project).First(&project).Error; err != nil {
		log.Println("Error Task: ", err)

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      fmt.Sprintf(types.ERR_MSG_NOT_FOUND, "project"),
		})
	}

	task.ProjectID = project.ID
	task.Project = project

	if taskErrFields, err := task.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err_fields": taskErrFields,
			"err_type":   types.ERR_TYPE_BY_MULTIPLE_FIELDS,
		})
	}

	if err := database.DB.Create(&task).Error; err != nil {
		log.Println("Error Task: ", err)

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{ // ANY UNEXPECTED ERR
			"msg":      fmt.Sprintf(types.ERR_MSG_SEVER_CREATE_ENTITY, "task"),
			"err_type": types.ERR_TYPE_MESSAGE,
		})
	}

	return c.Status(http.StatusOK).JSON(task.GetDictionary(isDefaultProject))
}
