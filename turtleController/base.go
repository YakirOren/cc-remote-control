package turtleController

import (
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Command struct {
	Action string
	Code   string
}

type TurtleInfo struct {
	ws        *melody.Session
	responses chan string
}

type TurtleController struct {
	Melody       *melody.Melody
	turtles      map[string]TurtleInfo
	turtleLock   *sync.Mutex
	responseLock *sync.Mutex
}

func New() *TurtleController {
	m := melody.New()
	controller := &TurtleController{
		Melody:     m,
		turtles:    make(map[string]TurtleInfo),
		turtleLock: new(sync.Mutex),
	}

	m.HandleConnect(controller.HandleConnect)
	m.HandleDisconnect(controller.HandleDisconnect)
	m.HandleMessage(controller.HandleMessage)

	return controller
}

func (ctrl *TurtleController) SendCommand(turtleID string, command Command) (string, error) {
	turtle, exists := ctrl.turtles[turtleID]
	if !exists {
		return "", fmt.Errorf("no turtle with ID '%s'", turtleID)
	}

	log.Infof("Sending commmand '%s'", command.Code)

	message, err := json.Marshal(command)
	if err != nil {
		return "", err
	}

	if err := turtle.ws.Write(message); err != nil {
		return "", err
	}

	return <-turtle.responses, nil

}
