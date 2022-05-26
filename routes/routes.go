package routes

import (
	"github.com/gofiber/fiber/v2"
	"samplegoapp.com/controller"
)

func Setup(app *fiber.App) {
	g := app.Group("/api")
	{
		g.Post("/register", controller.Register)
		g.Post("/login", controller.Login)
		// g.Get("/user", controller.GetUsers)
		g.Post("/logout", controller.Logout)

		u :=g.Group("/user")
		{
			u.Get("/", controller.GetUsers)

		p := u.Group("/post")
			{
				p.Get("/" , controller.GetPostsByUser)
				p.Post("/new", controller.MakePost)
				p.Get("/:id", controller.GetPostByUser)
				p.Put("/:id", controller.UpdatePost)
				p.Delete("/:id", controller.DeletePost)
				
			}
		}
	}

}
