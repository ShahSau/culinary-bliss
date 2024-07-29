package controllers

import (
	"fmt"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var restaurantCollection *mongo.Collection = database.GetCollection(database.DB, "restaurants")

func GetRestaurants(c *gin.Context) {
	restaurants, err := restaurantCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer restaurants.Close(c.Request.Context())

	var results []models.Restaurant

	for restaurants.Next(c.Request.Context()) {
		var restaurant models.Restaurant
		if err = restaurants.Decode(&restaurant); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		results = append(results, restaurant)
	}

	c.JSON(200, gin.H{"error": false, "message": "Restaurants retrived successfully", "data": results, "status": 200, "success": true})
}

func GetRestaurant(c *gin.Context) {
	restaurant_id := c.Param("id")

	if err := c.ShouldBindJSON(&restaurant_id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var restaurant models.Restaurant

	defer c.Request.Body.Close()

	err := restaurantCollection.FindOne(c.Request.Context(), bson.M{"restaurant_id": restaurant_id}).Decode(&restaurant)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"error": false, "message": "Restaurant retrived successfully", "data": restaurant, "status": 200, "success": true})
}

func CreateRestaurant(c *gin.Context) {
	fmt.Println("CreateRestaurant")
}

func UpdateRestaurant(c *gin.Context) {
	fmt.Println("UpdateRestaurant")
}

func DeleteRestaurant(c *gin.Context) {
	restaurant_id := c.Param("id")

	if err := c.ShouldBindJSON(&restaurant_id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := restaurantCollection.DeleteOne(c.Request.Context(), bson.M{"restaurant_id": restaurant_id})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"error": false, "message": "Restaurant deleted successfully", "status": 200, "success": true, "data": nil})
}
