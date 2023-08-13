package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func broadcast(m *melody.Melody, action string, userId string, msg string) {
	
		// Define the dynamic data
		data := struct {
			User    string
			Content string
			Color   string
			Time    string
		}{
			User:		userId,
			Content: msg,
			Color:   stringToHue(userId),
			Time:    time.Now().Format("15:04:05"),
		}

		// Parse the HTML template
		tmpl, err := template.ParseFiles("templates/"+ action +".html")
		if err != nil {
			log.Fatal("template parsing error: ", err)
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
}

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

		broadcast(m, "join", userId, "has joined the chat")
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
		
		broadcast(m,"message", userId, userMessage.Content)
	})

	r.GET("/", func(c *gin.Context) {
		// Render the "index.html" template, passing in a map of data
		c.HTML(200, "index.html", gin.H{})
	})

	r.Run(":3000")
}

func stringToHue(name string) string {
	hash := 0

	// Hash the name by summing the ASCII values of its characters
	for i := 0; i < len(name); i++ {
		hash += int(name[i])
	}

	// Calculate the hue value based on the hash
	// The hue value is in the range [0, 360], so we use the modulo operator
	hue := hash % 360
	return fmt.Sprintf("%d", hue)
}
