package controllers

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer restaurants.Close(c.Request.Context())

	var results []models.Restaurant

	for restaurants.Next(c.Request.Context()) {
		var restaurant models.Restaurant
		if err = restaurants.Decode(&restaurant); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		results = append(results, restaurant)
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Restaurants retrived successfully", "data": results, "status": http.StatusOK, "success": true})
}

func GetRestaurant(c *gin.Context) {
	restaurant_id := c.Param("id")

	if err := c.ShouldBindJSON(&restaurant_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var restaurant models.Restaurant

	defer c.Request.Body.Close()

	err := restaurantCollection.FindOne(c.Request.Context(), bson.M{"restaurant_id": restaurant_id}).Decode(&restaurant)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Restaurant retrived successfully", "data": restaurant, "status": http.StatusOK, "success": true})
}

func CreateRestaurant(c *gin.Context) {
	var restaurantReq models.Restaurant

	if err := c.ShouldBindJSON(&restaurantReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	var menuCollection *mongo.Collection = database.GetCollection(database.DB, "menus")

	var menus []models.Menu

	for _, item := range restaurantReq.Menu {
		cursor, err := menuCollection.Find(c.Request.Context(), bson.M{"menu_id": item.Menu_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(c.Request.Context())

		for cursor.Next(c.Request.Context()) {
			var menu models.Menu
			if err := cursor.Decode(&menu); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			menus = append(menus, menu)
		}
	}

	restaurantReq.Menu = menus

	_, err := restaurantCollection.InsertOne(c.Request.Context(), restaurantReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Restaurant created successfully", "status": http.StatusOK, "success": true, "data": restaurantReq})
}

func UpdateRestaurant(c *gin.Context) {
	restaurant_id := c.Param("id")

	if err := c.ShouldBindJSON(&restaurant_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var restaurant models.Restaurant

	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	var menuCollection *mongo.Collection = database.GetCollection(database.DB, "menus")

	var menus []models.Menu

	for _, item := range restaurant.Menu {
		cursor, err := menuCollection.Find(c.Request.Context(), bson.M{"menu_id": item.Menu_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(c.Request.Context())

		for cursor.Next(c.Request.Context()) {
			var menu models.Menu
			if err := cursor.Decode(&menu); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			menus = append(menus, menu)
		}
	}

	restaurant.Menu = menus

	_, err := restaurantCollection.UpdateOne(c.Request.Context(), bson.M{"restaurant_id": restaurant_id}, bson.M{"$set": restaurant})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Restaurant updated successfully", "status": http.StatusOK, "success": true, "data": restaurant})
}

func DeleteRestaurant(c *gin.Context) {
	restaurant_id := c.Param("id")

	if err := c.ShouldBindJSON(&restaurant_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	_, err := restaurantCollection.DeleteOne(c.Request.Context(), bson.M{"restaurant_id": restaurant_id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Restaurant deleted successfully", "status": http.StatusOK, "success": true, "data": nil})
}
