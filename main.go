package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"

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

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://127.0.0.1:3000",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
	}))
	initRoutes(app, inHandler, outHandler)

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}

}
