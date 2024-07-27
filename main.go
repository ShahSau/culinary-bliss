package main

import (
	"os"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/middleware"
	"github.com/ShahSau/culinary-bliss/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.GetCollection(database.DB, "food")

//var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
//var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")
//var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")
//var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")
//var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

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
