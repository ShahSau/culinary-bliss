package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func GlobalRoutes(c *gin.Engine) {
	c.GET("/categories", controllers.GetCategories)
	c.GET("/table", controllers.GetTables)
	c.GET("/table/:id", controllers.GetTable)
}
