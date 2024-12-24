package main

import (
	"os"

	"time"

	"github.com/ShahSau/culinary-bliss/database"
	docs "github.com/ShahSau/culinary-bliss/docs"
	"github.com/ShahSau/culinary-bliss/middleware"
	"github.com/ShahSau/culinary-bliss/routes"
	"github.com/gin-contrib/cors"
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
	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://culinary-bliss.onrender.com", "http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Swagger docs
	docs.SwaggerInfo.Title = "Culinary Bliss API"
	docs.SwaggerInfo.Description = "Culinary Bliss restaurant management app designed to streamline your operations and elevate your dining experience and for efficient, effective, and effortless restaurant management."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"https"}
	docs.SwaggerInfo.Host = "culinary-bliss.onrender.com"
	//docs.SwaggerInfo.Host = "localhost:8080"

	routes.AuthRoutes(router)
	routes.GlobalRoutes(router)
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
