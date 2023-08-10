package main

import (
	"log"
	"myContacts/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/handlebars/v2"
	_ "github.com/mattn/go-sqlite3"
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
        
        return c.Render("index", fiber.Map{
            "Title": "Hello, World!",
            "todos": todos,
                
        }, "layouts/main")
    })

    app.Get("/add", func (c *fiber.Ctx) error {
              
        return c.Render("add", fiber.Map{
        }, "layouts/main")
 
    })

    app.Post("/add", func (c* fiber.Ctx)error {
         todo :=new(db.TodoDraft)

        if err:= c.BodyParser(todo); err != nil {
            return err
        }

       db.AddTodo(todo) 
       return c.Redirect("/") 
    })

    log.Fatal(app.Listen(":3000"))
}
