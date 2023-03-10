package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Command struct {
	Action string
	Code   string
}

func (tc *TurtleController) SendCommand(turtleID string, command Command) (string, error) {
	turtle, exists := tc.turtles[turtleID]
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

// GetConnected godoc
// @Summary get active sessions
// @Schemes
// @Description get active sessions
// @Tags session
// @Accept json
// @Produce json
// @Success 200 {array} string "array of connected turtle IDs"
// @Router /api/v1/sessions [get]
func (tc *TurtleController) GetConnected(c *gin.Context) {
	c.JSON(http.StatusOK, Keys(tc.turtles))

}

// RunCommand godoc
// @Summary run command
// @Schemes
// @Description send command to turtle
// @Param id query string true "ID"
// @Param command body Command true "commmand to send" Command{} Command{Action: "eval", Code: "ls"}
// @Tags session
// @Accept json
// @Produce json
// @Success 200 {string} string "response from the turtle"
// @Router /api/v1/command [post]
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
