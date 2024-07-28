package controllers

import (
	"fmt"
	"net/http"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

func GetOrders(c *gin.Context) {
	orders, err := orderCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer orders.Close(c.Request.Context())

	var results []models.Order

	for orders.Next(c.Request.Context()) {
		var order models.Order
		if err = orders.Decode(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		}
		results = append(results, order)
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Order retrived successfully", "data": results, "status": http.StatusOK, "success": true})
}

func GetOrder(c *gin.Context) {
	fmt.Println("GetOrder")
}

func CreateOrder(c *gin.Context) {
	fmt.Println("CreateOrder")
}

func UpdateOrder(c *gin.Context) {
	fmt.Println("UpdateOrder")
}

func DeleteOrder(c *gin.Context) {
	fmt.Println("DeleteOrder")
}
