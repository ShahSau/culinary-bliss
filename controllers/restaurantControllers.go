package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, errp := strconv.Atoi(c.Query("page"))
	if errp != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(c.Query("startIndex"))

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "id", Value: 1},
				{Key: "title", Value: 1},
				{Key: "image", Value: 1},
				{Key: "time", Value: 1},
				{Key: "pickup", Value: 1},
				{Key: "delivery", Value: 1},
				{Key: "rating", Value: 1},
				{Key: "rating_count", Value: 1},
				{Key: "menu_id", Value: 1},
			},
		},
	}

	record := bson.D{{Key: "$skip", Value: recordPerPage * (page - 1)}}
	limit := bson.D{{Key: "$limit", Value: recordPerPage}}

	result, errAgg := restaurantCollection.Aggregate(c.Request.Context(), mongo.Pipeline{matchStage, projectStage, record, limit})

	if errAgg != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errAgg.Error()})
		return
	}

	var allRestaurants []bson.M
	if err = result.All(c.Request.Context(), &allRestaurants); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": allRestaurants, "page": page, "recordPerPage": recordPerPage, "startIndex": startIndex})
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

	var restaurant models.Restaurant

	defer c.Request.Body.Close()

	err := restaurantCollection.FindOne(c.Request.Context(), bson.M{"restaurant_id": restaurant_id}).Decode(&restaurant)

	if err != nil {
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
	var restaurantReq types.Restaurant

	if err := c.ShouldBindJSON(&restaurantReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	var restaurant models.Restaurant

	restaurant.ID = primitive.NewObjectID()
	restaurant.Restaurant_id = restaurant.ID.Hex()
	restaurant.Title = restaurantReq.Title
	restaurant.Image = restaurantReq.Image
	restaurant.Time = restaurantReq.Time
	restaurant.Pickup = restaurantReq.Pickup
	restaurant.Delivery = restaurantReq.Delivery
	restaurant.Rating = restaurantReq.Rating
	restaurant.RatingCount = restaurantReq.RatingCount
	restaurant.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	restaurant.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	var menus []models.Menu

	for _, item := range restaurantReq.Menu {
		var menu models.Menu
		err := menuCollection.FindOne(c.Request.Context(), bson.M{"menu_id": item}).Decode(&menu)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		menus = append(menus, menu)
	}
	restaurant.Menu = menus

	_, err := restaurantCollection.InsertOne(c.Request.Context(), restaurant)

	if err != nil {
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

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}

	var restaurantReq types.Restaurant

	if err := c.ShouldBindJSON(&restaurantReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	var restaurant models.Restaurant
	var menuCollection *mongo.Collection = database.GetCollection(database.DB, "menus")

	restaurant.Title = restaurantReq.Title
	restaurant.Image = restaurantReq.Image
	restaurant.Time = restaurantReq.Time
	restaurant.Pickup = restaurantReq.Pickup
	restaurant.Delivery = restaurantReq.Delivery
	restaurant.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	var menus []models.Menu

	for _, item := range restaurant.Menu {
		cursor, err := menuCollection.Find(c.Request.Context(), bson.M{"menu_id": item})
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

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
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
