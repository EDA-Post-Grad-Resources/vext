package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
	"log"
)

func main() {
	engine := handlebars.New("./views", ".hbs")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./public")

	app.Get("/", func(c *fiber.Ctx) error {

		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		}, "layouts/main")
	})

	log.Fatal(app.Listen(":3000"))
}
