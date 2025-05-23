package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func FoodRoutes(c *gin.Engine) {
	c.POST("/food", controllers.CreateFood)
	c.PUT("/food/:id", controllers.UpdateFood)
	c.DELETE("/food/:id", controllers.DeleteFood)
}
