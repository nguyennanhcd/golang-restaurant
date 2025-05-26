package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(menuRoutes *gin.Engine) {
	menuRoutes.GET("/menus", controller.GetMenus())
	menuRoutes.GET("/menus/:menu_id", controller.GetMenu())
	menuRoutes.POST("/menu", controller.CreateMenu())
	menuRoutes.PATCH("/menu/:menu_id", controller.UpdateMenu())
}
