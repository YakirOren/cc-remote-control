package main

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

type Command struct {
	Lua string
}

type TurtleInfo struct {
	Meta string
}

func main() {
	r := gin.Default()
	m := melody.New()
	turtles := make(map[string]*melody.Session)
	lock := new(sync.Mutex)

	r.POST("/command", runCommand(turtles))

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		lock.Lock()

		turtleID := s.Request.Header.Get("User-Agent")
		log.Infof("New turtle connected with ID: %s", turtleID)
		turtles[turtleID] = s

		lock.Unlock()
	})

	m.HandleDisconnect(func(s *melody.Session) {
		lock.Lock()

		turtleID := s.Request.Header.Get("User-Agent")
		log.Infof("turtle %s disconnected", turtleID)
		delete(turtles, turtleID)

		lock.Unlock()
	})

	r.Run(":4000")
}

func runCommand(turtles map[string]*melody.Session) func(c *gin.Context) {
	return func(c *gin.Context) {
		turtleID, _ := c.GetQuery("id")

		turtle, exists := turtles[turtleID]
		if !exists {
			c.String(http.StatusNotFound, "no turtle with ID '%s'", turtleID)
			return
		}

		var command Command

		if err := c.BindJSON(&command); err != nil {
			log.Error(err)
			return
		}

		log.Infof("running commmand '%s'", command.Lua)
		err := turtle.Write([]byte(command.Lua))
		if err != nil {
			log.Error(err)
			return
		}

	}
}
