package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(orderRoutes *gin.Engine) {
	orderRoutes.GET("/orders", controller.GetOrders())
	orderRoutes.GET("/orders/:order_id", controller.GetOrder())
	orderRoutes.POST("/order", controller.CreateOrder())
	orderRoutes.PATCH("/orders/:order_id", controller.UpdateOrder())
}
