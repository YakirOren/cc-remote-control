package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"net/http"
	"sync"
)

type TurtleController struct {
	Melody     *melody.Melody
	turtles    map[string]Session
	turtleLock *sync.RWMutex
}

type Session struct {
	ws        *melody.Session
	responses chan string
}

func New() *TurtleController {
	m := melody.New()
	controller := &TurtleController{
		Melody:     m,
		turtles:    make(map[string]Session),
		turtleLock: new(sync.RWMutex),
	}

	m.HandleConnect(controller.HandleConnect)
	m.HandleDisconnect(controller.HandleDisconnect)
	m.HandleMessage(controller.HandleMessage)

	return controller
}

func (tc *TurtleController) CreateWebsocket(c *gin.Context) {
	err := tc.Melody.HandleRequest(c.Writer, c.Request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
