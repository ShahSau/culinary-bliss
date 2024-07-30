package routes

import (
	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(c *gin.Engine) {
	c.POST("/login", controllers.Login)
	c.POST("/register", controllers.Register)
	c.POST("/logout", controllers.Logout)
}
