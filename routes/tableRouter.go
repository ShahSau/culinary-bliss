package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func TableRoutes(c *gin.Engine) {
	c.GET("/table", controllers.GetTables)
	c.GET("/table/:id", controllers.GetTable)
	c.POST("/table", controllers.CreateTable)
	c.PUT("/table/:id", controllers.UpdateTable)
	c.DELETE("/table/:id", controllers.DeleteTable)
}
