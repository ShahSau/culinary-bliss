package controllers

import (
	"log"
	"net/http"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/services"
	"github.com/ShahSau/culinary-bliss/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var restaurantCollection *mongo.Collection = database.GetCollection(database.DB, "restaurants")

// @Summary GetRestaurants
// @Description Get all restaurants
// @Tags Global
// @Produce json
// @Param 		 recordPerPage query int false "Record Per Page"
// @Param 		 page query int false "Page"
// @Param 		 startIndex query int false "Start Index"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /restaurants [get]
func GetRestaurants(c *gin.Context) {
	responseRestaurant, err := services.GetRestaurants(c)
	if err != nil {
		log.Println("Error getting restaurants:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": responseRestaurant.AllRestaurants, "page": responseRestaurant.Page, "recordPerPage": responseRestaurant.RecordPerPage, "startIndex": responseRestaurant.StartIndex, "error": false, "message": "Restaurants retrieved successfully", "status": http.StatusOK, "success": true})
}

// @Summary GetRestaurant
// @Description Get a restaurant
// @Tags Global
// @Produce json
// @Param id path string true "Restaurant ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /restaurants/{id} [get]
func GetRestaurant(c *gin.Context) {
	restaurant_id := c.Param("id")

	restaurant, err := services.GetRestaurantByID(c, restaurant_id)
	if err != nil {
		log.Println("Error getting restaurant:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Restaurant retrived successfully", "data": restaurant, "status": http.StatusOK, "success": true})
}

// @Summary CreateRestaurant
// @Description Create a restaurant
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param restaurant body types.Restaurant true "Restaurant Object"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /restaurants [post]
func CreateRestaurant(c *gin.Context) {
	var restaurantReq models.Restaurant

	if err := c.ShouldBindJSON(&restaurantReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant, err := services.CreateRestaurant(c, restaurantReq)
	if err != nil {
		log.Println("Error creating restaurant:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Restaurant created successfully", "status": http.StatusOK, "success": true, "data": restaurant})
}

// @Summary UpdateRestaurant
// @Description Update a restaurant
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Restaurant ID"
// @Param restaurant body types.Restaurant true "Restaurant Object"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /restaurants/{id} [put]
func UpdateRestaurant(c *gin.Context) {
	restaurant_id := c.Param("id")

	if err := c.ShouldBindJSON(&restaurant_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var restaurantReq models.Restaurant
	if err := c.ShouldBindJSON(&restaurantReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	restaurant, err := services.UpdateRestaurant(c, restaurant_id, restaurantReq)
	if err != nil {
		log.Println("Error updating restaurant:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Restaurant updated successfully", "status": http.StatusOK, "success": true, "data": restaurant})
}

// @Summary DeleteRestaurant
// @Description Delete a restaurant
// @Tags Admin
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Restaurant ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /restaurants/{id} [delete]
func DeleteRestaurant(c *gin.Context) {
	restaurant_id := c.Param("id")

	_, err := restaurantCollection.DeleteOne(c.Request.Context(), bson.M{"restaurant_id": restaurant_id})
	if err != nil {
		log.Println("Error deleting restaurant:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Restaurant deleted successfully", "status": http.StatusOK, "success": true, "data": nil})
}

// @Summary MenuByRestaurant
// @Description Get all menus by restaurant
// @Tags User
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Restaurant ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /restaurants/menus/{id} [get]
func MenuByRestaurant(c *gin.Context) {
	restaurant_id := c.Param("id")

	var restaurant models.Restaurant

	defer c.Request.Body.Close()

	err := restaurantCollection.FindOne(c.Request.Context(), bson.M{"restaurant_id": restaurant_id}).Decode(&restaurant)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Restaurant retrived successfully", "data": restaurant.Menu, "status": http.StatusOK, "success": true})
}

// @Summary AddRatingtoRestaurant
// @Description Add rating to a restaurant
// @Tags User
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Restaurant ID"
// @Param rating body types.Rating true "Rating Object"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /restaurants/rating/{id} [put]
func AddRatingtoRestaurant(c *gin.Context) {
	restaurant_id := c.Param("id")

	var restaurant models.Restaurant

	defer c.Request.Body.Close()

	err := restaurantCollection.FindOne(c.Request.Context(), bson.M{"restaurant_id": restaurant_id}).Decode(&restaurant)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var rating types.Rating

	if err := c.ShouldBindJSON(&rating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurant.Rating = (float64(restaurant.Rating)*float64(restaurant.RatingCount) + float64(rating.Rating)) / float64(restaurant.RatingCount+1)
	restaurant.RatingCount = restaurant.RatingCount + 1

	_, err = restaurantCollection.UpdateOne(c.Request.Context(), bson.M{"restaurant_id": restaurant_id}, bson.M{"$set": restaurant})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Rating added successfully", "status": http.StatusOK, "success": true, "data": restaurant})

}
