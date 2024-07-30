package main

import (
	"os"

	"github.com/ShahSau/culinary-bliss/database"
	docs "github.com/ShahSau/culinary-bliss/docs"
	"github.com/ShahSau/culinary-bliss/middleware"
	"github.com/ShahSau/culinary-bliss/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Swagger docs
	docs.SwaggerInfo.Title = "Culinary Bliss API"
	docs.SwaggerInfo.Description = "Culinary Bliss restaurant management app designed to streamline your operations and elevate your dining experience and for efficient, effective, and effortless restaurant management."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "https://culinary-bliss.onrender.com"

	routes.AuthRoutes(router)

	router.Use(middleware.Authtication)

	routes.UserRoutes(router)
	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.InvoiceRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.RestaurantRoutes(router)
	routes.CatgeoryRoutes(router)

	router.Run(":" + port)

}
