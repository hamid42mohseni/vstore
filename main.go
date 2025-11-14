package main

import (
	"vstore/backend/database"
	"vstore/backend/routes"
	"vstore/backend/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitRedis("localhost:6379", "", 0)
	database.Connect()

	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080")

}
