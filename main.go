package main

import (
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
	initRoutes(app, inHandler, outHandler)

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}

}
