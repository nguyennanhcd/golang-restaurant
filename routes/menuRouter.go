package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(menuRoutes *gin.Engine) {
	menuRoutes.GET("/menu", controller.GetMenus())
	menuRoutes.GET("/menu/:menu_id", controller.GetMenu())
	menuRoutes.POST("/menu", controller.CreateMenu())
	menuRoutes.PATCH("/menu/:menu_id", controller.UpdateMenu())
}
