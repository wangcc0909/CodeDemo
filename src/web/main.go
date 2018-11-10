package main

import "github.com/gin-gonic/gin"

func main() {
	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	app.Run(":8080")
}
