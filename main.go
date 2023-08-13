package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func broadcast(m *melody.Melody, sessions []*melody.Session, isSender bool, userId string, msg string) {

	// Define the dynamic data
	data := struct {
		User     string
		Content  string
		Color    string
		Time     string
		IsSender bool
	}{
		User:     userId,
		Content:  msg,
		Color:    stringToHue(userId),
		Time:     time.Now().Format("15:04:05"),
		IsSender: isSender,
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("templates/message.html")
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

	// loop through the sessions
	for _, s := range sessions {
		// get the user ID from the session
		userId, _ := s.MustGet("userId").(string)
		fmt.Println("Sending message to " + userId)
	}

	// Broadcast the rendered HTML
	m.BroadcastMultiple(renderedHTML.Bytes(), sessions)
}

func getOtherSessions(sessions []*melody.Session, session *melody.Session) []*melody.Session {
	var otherSessions []*melody.Session
	for _, s := range sessions {
		if s != session {
			otherSessions = append(otherSessions, s)
		}
	}
	return otherSessions
}

func getInitials(fullName string) string {
	names := strings.Split(fullName, " ")
	if len(names) < 2 {
		return "Invalid name"
	}
	initials := string(names[0]) + "-" + string(names[1])
	return strings.ToLower(initials)
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

	// create an empty slice of sessions
	var sessions []*melody.Session

	m.HandleConnect(func(s *melody.Session) {
		// generate a random user ID
		userId := gofakeit.Name()
		// set the user ID on the session
		s.Set("userId", userId)
		fmt.Println(s.Request.RemoteAddr + " has been assigned to " + userId)

		// add the session to the slice
		sessions = append(sessions, s)
		otherSessions := getOtherSessions(sessions, s)

		broadcast(m, otherSessions, false, userId, " has joined the chat")
		broadcast(m, []*melody.Session{s}, true, userId, "joined the chat")

		initials := getInitials(userId)
		m.BroadcastMultiple([]byte("<p class='flex-none font-bold' id='avatar'>"+initials+"@chat$</p> "), []*melody.Session{s})
	})

	m.HandleDisconnect(func(s *melody.Session) {
		// delete the session from the slice
		for i, session := range sessions {
			if session == s {
				sessions = append(sessions[:i], sessions[i+1:]...)
				break
			}
		}

		// Get the user ID from the session
		userId, exists := s.MustGet("userId").(string)
		if !exists {
			log.Fatal("User ID not found")
			return
		}

		broadcast(m, []*melody.Session{s}, true, userId, "left the chat")
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

		otherSessions := getOtherSessions(sessions, s)

		broadcast(m, []*melody.Session{s}, true, userId, userMessage.Content)
		broadcast(m, otherSessions, false, userId, userMessage.Content)
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
