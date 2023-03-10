package controller

import (
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
)

func (tc *TurtleController) HandleDisconnect(s *melody.Session) {
	turtleID := s.Request.Header.Get("User-Agent")

	tc.turtleLock.Lock()
	delete(tc.turtles, turtleID)
	tc.turtleLock.Unlock()

	log.Infof("turtle %s disconnected", turtleID)
}

func (tc *TurtleController) HandleConnect(s *melody.Session) {
	turtleID := s.Request.Header.Get("User-Agent")

	tc.turtleLock.Lock()
	tc.turtles[turtleID] = Session{
		ws:        s,
		responses: make(chan string, 1),
	}

	tc.turtleLock.Unlock()
	log.Infof("New turtle connected with ID: %s", turtleID)
}

func (tc *TurtleController) HandleMessage(s *melody.Session, msg []byte) {
	turtleID := s.Request.Header.Get("User-Agent")

	turtle := tc.turtles[turtleID]
	turtle.responses <- string(msg)
}
