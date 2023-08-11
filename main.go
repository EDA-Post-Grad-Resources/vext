package main

import (
	"bytes"
	"html/template"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	r := gin.Default()
	m := melody.New()

	r.Use(static.Serve("/public", static.LocalFile("./public", true)))

	// Load templates from the "templates" directory
	r.LoadHTMLGlob("templates/*")

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		// Define the dynamic data
		data := struct {
			Header  string
			Content string
		}{
			Header:  "Hello",
			Content: "This is the main page.",
		}

		// Parse the HTML template
		tmpl, err := template.ParseFiles("templates/message.html")
		if err != nil {
			// Handle error
			return
		}

		// Render the template with the data
		var renderedHTML bytes.Buffer
		err = tmpl.Execute(&renderedHTML, data)
		if err != nil {
			// Handle error
			return
		}

		// Broadcast the rendered HTML
		m.Broadcast(renderedHTML.Bytes())
	})

	r.GET("/", func(c *gin.Context) {
		// Render the "index.html" template, passing in a map of data
		c.HTML(200, "index.html", gin.H{
			"Title":   "Vext",
			"Header":  "Welcome!",
			"Content": "This is the main page.",
		})
	})

	r.Run(":3000")
}
