package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func CatgeoryRoutes(c *gin.Engine) {
	c.GET("/categeory/:id", controllers.GetCategoryByID)
	c.POST("/categeory", controllers.CreateCategory)
	c.PUT("/categeory/:id", controllers.UpdateCategory)
	c.DELETE("/categeory/:id", controllers.DeleteCategory)

}
