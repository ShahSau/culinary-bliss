package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func MenuRoutes(c *gin.Engine) {
	c.GET("/menu", controllers.GetMenus)
	c.GET("/menu/:id", controllers.GetMenu)
	c.POST("/menu", controllers.CreateMenu)
	c.PUT("/menu/:id", controllers.UpdateMenu)
	c.DELETE("/menu/:id", controllers.DeleteMenu)

}
