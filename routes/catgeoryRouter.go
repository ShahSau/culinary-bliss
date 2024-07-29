package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func CatgeoryRoutes(c *gin.Engine) {
	c.GET("/categories", controllers.GetRestaurants)
	c.GET("/categeory/:id", controllers.GetRestaurant)
	c.POST("/categories", controllers.CreateRestaurant)
	c.PUT("/categeory/:id", controllers.UpdateRestaurant)
	c.DELETE("/categeory/:id", controllers.DeleteRestaurant)

}
