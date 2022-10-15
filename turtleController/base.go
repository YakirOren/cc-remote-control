package turtleController

import (
	"fmt"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Command struct {
	Lua string
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

func (ctrl *TurtleController) RunCommand(turtleID string, command Command) (string, error) {
	turtle, exists := ctrl.turtles[turtleID]
	if !exists {
		return "", fmt.Errorf("no turtle with ID '%s'", turtleID)
	}

	log.Infof("running commmand '%s'", command.Lua)
	err := turtle.ws.Write([]byte(command.Lua))
	if err != nil {
		return "", err
	}

	return <-turtle.responses, nil

}
