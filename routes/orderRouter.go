package routes

import (
	"fmt"

	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func OrderRoutes(c *gin.Engine) {
	fmt.Println("OrderRoutes")
	c.GET("/order", controllers.GetOrders)
	c.GET("/order/:id", controllers.GetOrder)
	c.POST("/order", controllers.CreateOrder)
	c.PUT("/order/:id", controllers.UpdateOrder)
	c.DELETE("/order/:id", controllers.DeleteOrder)
}
