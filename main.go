package main

import (
	"log"
	"myContacts/db"
	"strconv"

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

    app.Post("/delete/:id", func (c *fiber.Ctx)error {
        inputId := c.Params("id")

        // convert string to int
        id, err := strconv.Atoi(inputId)
        if err != nil {
            return err
        }

        db.Delete(id)
        return c.Redirect("/") 
    })

    app.Get("/edit/:id", func (c *fiber.Ctx)error {
        inputId := c.Params("id")

        // convert string to int
        id, err := strconv.Atoi(inputId)
        if err != nil {
            return err
        }

        todo := db.GetTodoById(id)
        log.Println(todo)
        return c.Render("edit", fiber.Map{
            "todo": db.Todo{
                Id: todo.Id,
                Task:todo.Task,
            },
        }, "layouts/main")
    })

    app.Post("/edit/:id", func (c *fiber.Ctx)error {
        inputId := c.Params("id")

        // convert string to int
        id, err := strconv.Atoi(inputId)
        if err != nil {
            return err
        }

        todo :=new(db.TodoDraft)

        if err:= c.BodyParser(todo); err != nil {
            return err
        }

        db.UpdateTodo(id, todo)
        return c.Redirect("/")
    })


    log.Fatal(app.Listen(":3000"))
}
