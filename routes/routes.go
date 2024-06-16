package routes

import (
	"github.com/gofiber/fiber/v2"
	"todolistserver.com/test/authenticator"
	"todolistserver.com/test/controllers"
	"todolistserver.com/test/midlewares"
)

func registerProjectRoutes(api fiber.Router) {
	projectsRoute := api.Group("/projects")

	projectsRoute.Get("/", controllers.GetAllProjects)
	projectsRoute.Get("/:id", controllers.GetProjectById)
	projectsRoute.Post("/", controllers.RegisterProject)
	projectsRoute.Put("/:id", controllers.UpdateProject)
	projectsRoute.Delete("/:id", controllers.DeleteProject)

	projectsRoute.Get("/:projectid/tasks", controllers.GetAllTasksByProjectId)
}

func registerProjectTasksRoutes(api fiber.Router) {
	projectsRoute := api.Group("/tasks")

	projectsRoute.Get("/", controllers.GetAllTasks)
	projectsRoute.Post("/", controllers.RegisterTask)
	projectsRoute.Put("/:id", controllers.UpdateTask)
	projectsRoute.Put("/:id/star", controllers.UpdateStarredTask)
	projectsRoute.Delete("/:id", controllers.DeleteTask)
}

func Register(app *fiber.App, auth *authenticator.Authenticator) {
	api := app.Group("/api")

	api.Use(midlewares.RouteMilewareAuth(auth))
	api.Use(midlewares.PrepeareProjectDefault)

	//store := session.New()

	//api.Use(session.Session("auth-session", store))

	// SET GROUPS -------------------------
	registerProjectRoutes(api)
	registerProjectTasksRoutes(api)
}
