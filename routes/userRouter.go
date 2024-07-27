package routes

import (
	"fmt"

	"github.com/ShahSau/culinary-bliss/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(c *gin.Engine) {
	c.GET("/users", controllers.GetUsers)
	c.GET("/users/:id", controllers.GetUser)
	c.POST("/users", controllers.CreateUser)
	c.PUT("/users/:id", controllers.UpdateUser)
	c.DELETE("/users/:id", controllers.DeleteUser)
	c.POST("/login", controllers.Login)
	c.POST("/register", controllers.Register)
	fmt.Println("UserRoutes")
}
