package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func TableRoutes(tableRoutes *gin.Engine) {
	tableRoutes.GET("/tables", controller.GetTables())
	tableRoutes.GET("/tables/:table_id", controller.GetTable())
	tableRoutes.POST("/tables", controller.CreateTable())
	tableRoutes.PATCH("/tables/:table_id", controller.UpdateTable())
}
