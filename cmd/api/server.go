package main

import (
	"turtle/controller"
	_ "turtle/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	app.Run(":8080")
}

func handleSwagger(app *gin.Engine) {
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
