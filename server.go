package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	log "github.com/sirupsen/logrus"
	"net/http"
	"turtle/controller"
	"turtle/docs"
	_ "turtle/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			Turtle C2 API
// @version		1.0
// @description	This server controls computer craft turtles.
// @BasePath	/
func main() {
	docs.SwaggerInfo.Host = ""

	app := gin.Default()
	app.Use(cors.Default())
	app.Use(gzip.Gzip(gzip.DefaultCompression))

	app.Use(static.Serve("/", static.LocalFile("static", false)))
	app.NoRoute(func(c *gin.Context) {
		c.File("static/index.html")
	})

	tc := controller.New()

	app.GET("/ws", tc.CreateWebsocket)

	v1group := app.Group("/api/v1")

	v1group.GET("/sessions", tc.GetConnected)
	v1group.POST("/command", tc.RunCommand)
	v1group.DELETE("/disconnect", tc.Disconnect)

	handleSwagger(app)

	log.Error(app.Run(":8080"))
}

func handleSwagger(app *gin.Engine) {
	app.GET("/swagger/:any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.GET("/swagger", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}
