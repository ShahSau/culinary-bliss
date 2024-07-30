package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(c *gin.Engine) {
	c.GET("/food", controllers.GetFoods)
	c.GET("/food/:id", controllers.GetFood)
	c.POST("/food", controllers.CreateFood)
	c.PUT("/food/:id", controllers.UpdateFood)
	c.DELETE("/food/:id", controllers.DeleteFood)
}
