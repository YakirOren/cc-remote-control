package turtleController

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"net/http"
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

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
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

func (tc *TurtleController) GetConnected(c *gin.Context) {
	c.JSON(http.StatusOK, Keys(tc.turtles))
}

func (tc *TurtleController) RunCommand(c *gin.Context) {

	turtleID, _ := c.GetQuery("id")

	var command Command
	if err := c.BindJSON(&command); err != nil {
		log.Error(err)
		c.String(http.StatusBadRequest, "cant parse request body")
		return
	}

	response, err := tc.SendCommand(turtleID, command)
	if err != nil {
		log.Error(err)
		c.String(http.StatusInternalServerError, "running command '%s' failed: %w", command, err)
		return
	}

	c.String(http.StatusOK, response)

}

func (tc *TurtleController) CreateWebsocket(c *gin.Context) {
	tc.Melody.HandleRequest(c.Writer, c.Request)
}
