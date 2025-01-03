package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func OrderItemRoutes(c *gin.Engine) {
	c.GET("/orderItem", controllers.GetOrderItems)
	c.GET("/orderItem/:id", controllers.GetOrderItem)
	c.POST("/orderItem", controllers.CreateOrderItem)
	c.PUT("/orderItem/:id", controllers.UpdateOrderItem)
	c.DELETE("/orderItem/:id", controllers.DeleteOrderItem)
}
