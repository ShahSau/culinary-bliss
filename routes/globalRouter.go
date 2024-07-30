package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func GlobalRoutes(c *gin.Engine) {
	c.GET("/categories", controllers.GetCategories)
	c.GET("/table", controllers.GetTables)
	c.GET("/table/:id", controllers.GetTable)
	c.GET("/menu", controllers.GetMenus)
	c.GET("/menu/:id", controllers.GetMenu)
	c.GET("/restaurants", controllers.GetRestaurants)
	c.GET("/restaurants/:id", controllers.GetRestaurant)
	c.GET("/restaurants/menus/:id", controllers.MenuByRestaurant)
	c.GET("/food", controllers.GetFoods)
	c.GET("/food/:id", controllers.GetFood)
}
