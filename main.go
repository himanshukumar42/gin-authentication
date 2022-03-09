package main

import (
	"github.com/gin-gonic/gin"
	"github.com/himanshuk42/gin-authentication/configs"
	"github.com/himanshuk42/gin-authentication/routes"
)

func main() {
	router := gin.Default()

	configs.ConnectDB()

	routes.UserRoute(router)

	router.Run("localhost:8090")
}
