package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
	_ "github.com/mattn/go-sqlite3" 
	"myContacts/db"
)

func main() {
    engine := handlebars.New("./views", ".hbs")
    app := fiber.New(fiber.Config{
        Views: engine,
    })
	app.Static("/", "./public")
    
    
    app.Get("/", func(c *fiber.Ctx) error {
        // Render index
        todos := db.GetTodos() 
       log.Println(todos) 
        return c.Render("index", fiber.Map{
            "Title": "Hello, World!",
            "todos": todos,
                
        }, "layouts/main")
    })

    log.Fatal(app.Listen(":3000"))
}
