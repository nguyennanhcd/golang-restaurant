package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(orderItemRoutes *gin.Engine) {
	orderItemRoutes.GET("/orderItems", controller.GetOrderItems())
	orderItemRoutes.GET("/orderItems/:orderItem_id", controller.GetOrderItem())
	orderItemRoutes.GET("/orderItems-order/:order_id", controller.GetOrderItemsByOrder())
	orderItemRoutes.POST("/orderItem", controller.CreateOrderItem())
	orderItemRoutes.PATCH("/orderItems/:orderItem_id", controller.UpdateOrderItem())
}
