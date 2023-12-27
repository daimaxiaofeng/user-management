package router

import (
	"github.com/daimaxiaofeng/user-management/handlers"
	"github.com/daimaxiaofeng/user-management/middlewares"
	"github.com/gin-gonic/gin"
)

type Route struct {
	RelativePath string
	HandlerFunc  gin.HandlerFunc
}

func Init() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(middlewares.Cors)

	r.POST("/register", handlers.RegisterHandler)

	return r
}
