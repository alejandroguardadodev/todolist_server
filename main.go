package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"todolistserver.com/test/database"
	"todolistserver.com/test/models"
	"todolistserver.com/test/routes"
	"todolistserver.com/test/validation"
)

func migrate() {
	log.Println("Running Migation")
	database.DB.AutoMigrate(&models.Project{})
}

func main() {
	godotenv.Load()

	database.Conect()

	migrate()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	url := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))

	validation.ValidationInit()

	routes.Register(app)

	log.Fatal(app.Listen(url))
}