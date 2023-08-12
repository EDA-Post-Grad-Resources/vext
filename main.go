package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"

	"github.com/brianvoe/gofakeit/v6"
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

	m.HandleConnect(func(s *melody.Session) {
		// generate a random user ID
		userId := gofakeit.Name()
		// set the user ID on the session
		s.Set("userId", userId)
		fmt.Println(s.Request.RemoteAddr + " has been assigned to " + userId)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {

		// parse the message into an object that has userMessage property
		var userMessage struct {
			Content string `json:"userMessage"`
		}
		// unmarshal the message into the object
		err := json.Unmarshal(msg, &userMessage)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Get the user ID from the session
		userId, exists := s.MustGet("userId").(string)
		if !exists {
			log.Fatal("User ID not found")
			return
		}
		fmt.Println("User ID: ", userMessage.Content, userId)

		// Define the dynamic data
		data := struct {
			User    string
			Content string
		}{
			User:    userId,
			Content: userMessage.Content,
		}

		// Parse the HTML template
		tmpl, err := template.ParseFiles("templates/message.html")
		if err != nil {
			log.Fatal(err)
			return
		}

		// Render the template with the data
		var renderedHTML bytes.Buffer
		err = tmpl.Execute(&renderedHTML, data)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Broadcast the rendered HTML
		m.Broadcast(renderedHTML.Bytes())
	})

	r.GET("/", func(c *gin.Context) {
		// Render the "index.html" template, passing in a map of data
		c.HTML(200, "index.html", gin.H{})
	})

	r.Run(":3000")
}
