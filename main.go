package main

import (
	// static "github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	r := gin.Default()
	m := melody.New()

	// r.Use(static.Serve("/", static.LocalFile("./public", true)))

	// Load templates from the "templates" directory
	r.LoadHTMLGlob("templates/*")

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		htmlString := `<h1>Hello</h1>`
		m.Broadcast([]byte(htmlString))
	})

	r.GET("/", func(c *gin.Context) {
		// Render the "index.html" template, passing in a map of data
		c.HTML(200, "index.html", gin.H{
			"Title":   "My Gin Website",
			"Header":  "Welcome!",
			"Content": "This is the main page.",
		})
	})

	r.Run(":3000")
}
