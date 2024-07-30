package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func RestaurantRoutes(c *gin.Engine) {
	c.POST("/restaurants", controllers.CreateRestaurant)
	c.PUT("/restaurants/:id", controllers.UpdateRestaurant)
	c.DELETE("/restaurants/:id", controllers.DeleteRestaurant)
	c.PUT("/restaurants/rating/:id", controllers.AddRatingtoRestaurant)

}
