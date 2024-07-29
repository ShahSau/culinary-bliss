package main

import (
	"os"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/middleware"
	"github.com/ShahSau/culinary-bliss/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.ConnectDB()

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	routes.UserRoutes(router)
	router.Use(middleware.Authtication)

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.InvoiceRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)

	router.Run(":" + port)

}
