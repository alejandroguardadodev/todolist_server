package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"todolistserver.com/test/database"
	"todolistserver.com/test/models"
	"todolistserver.com/test/types"
)

func GetAllTasks(c *fiber.Ctx) error {
	user := c.Locals("user").(string)

	limit, limiterr := strconv.Atoi(c.Query("limit", "10"))
	page, pageerr := strconv.Atoi(c.Query("page", "0"))
	orderby := c.Query("orderby", "created_at")
	order := c.Query("order", "asc")

	if limiterr != nil {
		limit = 10
	}

	if pageerr != nil {
		page = 0
	}

	var projectsID []uint

	database.DB.Model(models.Project{}).Where(&models.Project{User: user}).Pluck("id", &projectsID)

	tasks := []models.Task{}

	if err := database.DB.Preload("Project").Where("project_id IN ?", projectsID).Order(fmt.Sprintf("%s %s", orderby, strings.ToUpper(order))).Limit(limit).Offset(page * limit).Find(&tasks).Error; err != nil {
		log.Println("Error Tasks: ", err)

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"err_type": types.ERR_TYPE_MESSAGE,
			"msg":      fmt.Sprintf(types.ERR_MSG_NOT_FOUND, "tasks"),
		})
	}

	var counts int64
	database.DB.Model(&models.Task{}).Where("project_id IN ?", projectsID).Count(&counts)

	_tasks := []models.Dictionary{}

	for _, task := range tasks {
		_tasks = append(_tasks, *task.GetDictionary(task.Project.Title == *models.GetDefaultProjectTitle(user)))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"tasks": _tasks,
		"total": counts,
	})
}

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

	task.AdjustDates()

	if taskErrFields, err := task.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"fields":   taskErrFields,
			"err_type": types.ERR_TYPE_BY_MULTIPLE_FIELDS,
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
