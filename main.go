package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"

	"github.com/gofiber/fiber/v2"

	"main/handlers"
	"main/repository"
	"main/services"
)

func main() {
	repos := repository.InitRepo(os.Getenv("DB_CONNECTION"))
	repos.Migrate()

	user := services.NewUsers(repos)

	outHandler := handlers.NewOutHandlers(repos, user)
	inHandler := handlers.NewInHandlers(repos, user)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:5500,http://localhost:8080,http://host.docker.internal:8080,http://localhost:5500,http://127.0.0.1:3000,http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	initRoutes(app, inHandler, outHandler)

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}

}
