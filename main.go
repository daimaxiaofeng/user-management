package main

import (
	"github.com/daimaxiaofeng/user-management/router"
	"github.com/daimaxiaofeng/user-management/utils"
)

func main() {
	r := router.Init()
	defer utils.DB.Close()

	// Listen and Server in 0.0.0.0:8080
	r.Run(":2024")
}
