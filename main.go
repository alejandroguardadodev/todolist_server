package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"todolistserver.com/test/authenticator"
	"todolistserver.com/test/database"
	"todolistserver.com/test/models"
	"todolistserver.com/test/routes"
	"todolistserver.com/test/validation"
)

// The `migrate` function in Go logs a message and performs auto migration for the Project and Task
// models in the database.
func migrate() {
	log.Println("Running Migation")
	database.DB.AutoMigrate(&models.Project{}, &models.Task{})
}

func main() {
	godotenv.Load()

	database.Conect()

	migrate()

	app := fiber.New()

	auth, err := authenticator.New()

	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://moonnuittodolist.netlify.app",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	url := fmt.Sprintf(":%s", os.Getenv("PORT"))

	validation.ValidationInit()

	routes.Register(app, auth)

	log.Fatal(app.Listen(url))
}
