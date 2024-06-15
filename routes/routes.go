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
}

func registerProjectTasksRoutes(api fiber.Router) {
	projectsRoute := api.Group("/tasks")

	projectsRoute.Post("/", controllers.RegisterTask)
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
