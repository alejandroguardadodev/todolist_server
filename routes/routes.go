package routes

import (
	"github.com/gofiber/fiber/v2"
	"todolistserver.com/test/controllers"
)

func registerProjectRoutes(api fiber.Router) {
	projectsRoute := api.Group("/projects")

	projectsRoute.Get("/", controllers.GetAllProjects)
	projectsRoute.Get("/:id", controllers.GetProjectById)
	projectsRoute.Post("/", controllers.RegisterProject)
	projectsRoute.Put("/:id", controllers.UpdateProject)
}

func Register(app *fiber.App) {
	api := app.Group("/api")

	// SET GROUPS -------------------------
	registerProjectRoutes(api)
}
