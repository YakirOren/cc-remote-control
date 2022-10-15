package turtleController

import (
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
)

func (ctrl *TurtleController) HandleDisconnect(s *melody.Session) {
	ctrl.turtleLock.Lock()

	turtleID := s.Request.Header.Get("User-Agent")
	log.Infof("turtle %s disconnected", turtleID)
	delete(ctrl.turtles, turtleID)

	ctrl.turtleLock.Unlock()
}

func (ctrl *TurtleController) HandleConnect(s *melody.Session) {
	ctrl.turtleLock.Lock()

	turtleID := s.Request.Header.Get("User-Agent")
	log.Infof("New turtle connected with ID: %s", turtleID)
	ctrl.turtles[turtleID] = TurtleInfo{
		ws:        s,
		responses: make(chan string, 1),
	}

	ctrl.turtleLock.Unlock()
}

func (ctrl *TurtleController) HandleMessage(s *melody.Session, msg []byte) {
	turtleID := s.Request.Header.Get("User-Agent")
	turtle := ctrl.turtles[turtleID]

	turtle.responses <- string(msg)
}
