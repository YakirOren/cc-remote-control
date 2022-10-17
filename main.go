package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"turtle/turtleController"
)

func main() {
	r := gin.Default()

	tc := turtleController.New()

	r.POST("/command", RunCommand(tc))
	r.GET("/ws", CreateWebsocket(tc))

	r.Run(":4000")
}

func RunCommand(tc *turtleController.TurtleController) func(c *gin.Context) {
	return func(c *gin.Context) {
		turtleID, _ := c.GetQuery("id")

		var command turtleController.Command
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
}

func CreateWebsocket(tc *turtleController.TurtleController) func(c *gin.Context) {
	return func(c *gin.Context) {
		tc.Melody.HandleRequest(c.Writer, c.Request)
	}
}
