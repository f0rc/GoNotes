package routes

import (
	"github.com/gofiber/fiber/v2"
	"samplegoapp.com/controller"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)
	app.Get("/api/user", controller.GetUsers)
	app.Post("/api/logout", controller.Logout)
}
