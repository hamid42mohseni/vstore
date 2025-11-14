package routes

import (
	"vstore/backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/user/login", controllers.CreateUserOrLogin)
		api.POST("/user/sendcode", controllers.SendCode)
	}
}
