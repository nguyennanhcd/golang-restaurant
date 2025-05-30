package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.POST("/user/signup", controller.SignUp())
	incomingRoutes.POST("/user/login", controller.Login())
	incomingRoutes.POST("/user/refresh-token", controller.RefreshToken())
}
