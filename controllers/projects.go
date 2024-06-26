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

func GetAllProjects(c *fiber.Ctx) error {
	user := c.Locals("user").(string)
	projects := []models.Project{}

	if err := database.DB.Where(models.Project{User: user}).Order("created_at").Find(&projects).Error; err != nil {
		log.Println("Error Project: ", err)

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      fmt.Sprintf(types.ERR_MSG_NOT_FOUND, "projects"),
		})
	}

	_projects := []models.Dictionary{}

	for _, project := range projects {
		if project.Title != *models.GetDefaultProjectTitle(user) {
			_projects = append(_projects, *project.GetDictionary())
		}
	}

	return c.Status(http.StatusOK).JSON(_projects)
}

func GetProjectById(c *fiber.Ctx) error {
	user := c.Locals("user").(string)
	id, err := c.ParamsInt("id")

	if err != nil {
		log.Println("Error Project: Bad Request")

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      "Invalid project ID",
		})
	}

	project := models.Project{
		ID:   uint(id),
		User: user,
	}

	if err := database.DB.Where(project).First(&project).Error; err != nil {
		log.Println("Error Project: ", err)

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      fmt.Sprintf(types.ERR_MSG_NOT_FOUND, "project"),
		})
	}

	if project.Title == *models.GetDefaultProjectTitle(user) {
		log.Println("Error Project: ", "DEFAULT PROJECT GET BY ID ERR")

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      fmt.Sprintf(types.ERR_MSG_NOT_FOUND, "project"),
		})
	}

	return c.Status(http.StatusOK).JSON(project.GetDictionary())
}

func UpdateProject(c *fiber.Ctx) error {
	user := c.Locals("user").(string)

	id, err := c.ParamsInt("id")

	if err != nil {
		log.Println("Error Project: Bad Request")

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      "Invalid project ID",
		})
	}

	project := models.Project{
		ID:   uint(id),
		User: user,
	}

	if err := database.DB.Where(project).First(&project).Error; err != nil {
		log.Println("Error Project: ", err)

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      fmt.Sprintf(types.ERR_MSG_NOT_FOUND, "project"),
		})
	}

	if err := c.BodyParser(&project); err != nil {
		return c.Status(http.StatusBadRequest).SendString(types.ERR_MSG_BAR_BODY_PARSE)
	}

	if project.Title == *models.GetDefaultProjectTitle(user) {
		log.Println("Error Project: ", "DEFAULT PROJECT UPDATE ERR")

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      fmt.Sprintf(types.ERR_MSG_NOT_FOUND, "project"),
		})
	}

	if projectErrFields, err := project.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err_fields": projectErrFields,
			"err_type":   types.ERR_TYPE_BY_MULTIPLE_FIELDS,
		})
	}

	if err := database.DB.Model(&project).Updates(project).Error; err != nil {
		log.Println("Error Project: ", err)

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{ // ANY UNEXPECTED ERR
			"msg":      fmt.Sprintf(types.ERR_MSG_SEVER_UPDATE_ENTITY, "project"),
			"err_type": types.ERR_TYPE_MESSAGE,
		})
	}

	return c.Status(http.StatusOK).JSON(project.GetDictionary())
}

func RegisterProject(c *fiber.Ctx) error {
	user := c.Locals("user").(string)

	var project models.Project

	if err := c.BodyParser(&project); err != nil {
		return c.Status(http.StatusBadRequest).SendString(types.ERR_MSG_BAR_BODY_PARSE)
	}

	project.User = user

	if projectErrFields, err := project.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"err_fields": projectErrFields,
			"err_type":   types.ERR_TYPE_BY_MULTIPLE_FIELDS,
		})
	}

	if project.Title == *models.GetDefaultProjectTitle(user) {
		log.Println("Error Project: ", "DEFAULT PROJECT REGISTER ERR")

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      types.ERR_UNEXPECTED,
		})
	}

	var counts int64
	database.DB.Model(&models.Project{}).Where(project).Count(&counts)

	if counts > 0 {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_BY_MULTIPLE_FIELDS,
			"fields": map[string]string{
				"title": "This project already exists",
			},
		})
	}

	if err := database.DB.Create(&project).Error; err != nil {
		log.Println("Error Project: ", err)

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{ // ANY UNEXPECTED ERR
			"msg":      fmt.Sprintf(types.ERR_MSG_SEVER_CREATE_ENTITY, "project"),
			"err_type": types.ERR_TYPE_MESSAGE,
		})
	}

	return c.Status(http.StatusOK).JSON(project.GetDictionary())
}

func DeleteProject(c *fiber.Ctx) error {
	user := c.Locals("user").(string)

	id, err := c.ParamsInt("id")

	if err != nil {
		log.Println("Error Project: Bad Request")

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      "Invalid project ID",
		})
	}

	project := models.Project{
		ID:   uint(id),
		User: user,
	}

	var counts int64
	database.DB.Model(&models.Project{}).Where(models.Project{ID: uint(id), Title: *models.GetDefaultProjectTitle(user)}).Count(&counts)

	if counts > 0 {
		log.Println("Error Project: ", "DEFAULT PROJECT DELETE ERR")

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      types.ERR_UNEXPECTED,
		})
	}

	if err := database.DB.Where(&models.Task{ProjectID: project.ID}).Delete(&[]models.Task{}).Error; err != nil {
		log.Println("Error Tasts By Project: ", err)

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      types.ERR_UNEXPECTED,
		})
	}

	if err := database.DB.Where(project).Delete(&project).Error; err != nil {
		log.Println("Error Project: ", err)

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      fmt.Sprintf(types.ERR_MSG_NOT_FOUND, "project"),
		})
	}

	return c.Status(http.StatusOK).JSON(project.GetDictionary())
}
