package main

import (
	"github.com/gin-gonic/gin"
	"turtle/turtleController"
)

func main() {
	r := gin.Default()

	tc := turtleController.New()

	r.POST("/command", tc.RunCommand)
	r.GET("/ws", tc.CreateWebsocket)
	r.GET("/sessions", tc.GetConnected)

	r.Run(":4000")
}
