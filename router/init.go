package router

import (
	"github.com/daimaxiaofeng/user-management/handlers"
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

	r.POST("/register", handlers.RegisterHandler)

	return r
}
