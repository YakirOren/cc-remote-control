package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"turtle/controller"
	_ "turtle/docs"
)

// @title Turtle C2 API
// @version 1.0
// @description This server controls computer craft turtles.

// @host localhost:4000
// @BasePath /
func main() {
	app := gin.Default()
	tc := controller.New()

	app.GET("/ws", tc.CreateWebsocket)

	v1group := app.Group("/api/v1")

	v1group.GET("/sessions", tc.GetConnected)
	v1group.POST("/command", tc.RunCommand)

	handleSwagger(app)

	app.Run(":4000")
}

func handleSwagger(app *gin.Engine) {
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
