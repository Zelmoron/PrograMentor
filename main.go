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
		AllowOrigins:     "http://26.199.2.128:3000",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "SKFJBhjfk",
	}))
	app.Use(logger.New())
	initRoutes(app, inHandler, outHandler)

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}

}
