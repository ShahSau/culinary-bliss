package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var restaurantCollection *mongo.Collection = database.GetCollection(database.DB, "restaurants")

func GetRestaurants(c *gin.Context) (models.ResponseRestaurant, error) {
	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage

	matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "id", Value: 1},
		{Key: "name", Value: 1},
		{Key: "image", Value: 1},
		{Key: "description", Value: 1},
	}}}
	skipStage := bson.D{{Key: "$skip", Value: startIndex}}
	limitStage := bson.D{{Key: "$limit", Value: recordPerPage}}

	cursor, err := restaurantCollection.Aggregate(c.Request.Context(), mongo.Pipeline{
		matchStage, projectStage, skipStage, limitStage,
	})
	if err != nil {
		return models.ResponseRestaurant{}, err
	}
	defer cursor.Close(c.Request.Context())

	var restaurants []models.Restaurant
	for cursor.Next(c.Request.Context()) {
		var restaurant models.Restaurant
		if err := cursor.Decode(&restaurant); err != nil {
			return models.ResponseRestaurant{}, err
		}
		restaurants = append(restaurants, restaurant)
	}

	response := models.ResponseRestaurant{
		AllRestaurants: restaurants,
		Page:           page,
		RecordPerPage:  recordPerPage,
		StartIndex:     startIndex,
	}

	return response, nil
}

func GetRestaurantByID(c *gin.Context, id string) (models.Restaurant, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Restaurant{}, err
	}
	var restaurant models.Restaurant
	err = restaurantCollection.FindOne(c.Request.Context(), bson.M{"_id": objectID}).Decode(&restaurant)
	if err != nil {
		return models.Restaurant{}, err
	}

	return restaurant, nil
}

func CreateRestaurant(c *gin.Context, restaurantReq models.Restaurant) (models.Restaurant, error) {
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
			return models.Restaurant{}, err
		}
		menus = append(menus, menu)
	}
	restaurant.Menu = menus

	_, err := restaurantCollection.InsertOne(c.Request.Context(), restaurant)

	if err != nil {
		return models.Restaurant{}, err
	}

	return restaurant, nil
}

func UpdateRestaurant(c *gin.Context, id string, restaurantReq models.Restaurant) (models.Restaurant, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Restaurant{}, err
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.Restaurant{}, errors.New("unauthorized")
	}

	var restaurant models.Restaurant
	var menuCollection *mongo.Collection = database.GetCollection(database.DB, "menus")
	err = restaurantCollection.FindOne(c.Request.Context(), bson.M{"_id": objectID}).Decode(&restaurant)
	if err != nil {
		return models.Restaurant{}, err
	}

	restaurant.Title = restaurantReq.Title
	restaurant.Image = restaurantReq.Image
	restaurant.Time = restaurantReq.Time
	restaurant.Pickup = restaurantReq.Pickup
	restaurant.Delivery = restaurantReq.Delivery
	restaurant.Rating = restaurantReq.Rating
	restaurant.RatingCount = restaurantReq.RatingCount
	restaurant.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	var menus []models.Menu

	for _, item := range restaurantReq.Menu {
		var menu models.Menu
		err := menuCollection.FindOne(c.Request.Context(), bson.M{"menu_id": item}).Decode(&menu)
		if err != nil {
			return models.Restaurant{}, err
		}
		menus = append(menus, menu)
	}
	restaurant.Menu = menus

	update := bson.M{
		"$set": restaurant,
	}

	err = restaurantCollection.FindOneAndUpdate(c.Request.Context(), bson.M{"_id": objectID}, update).Decode(&restaurant)
	if err != nil {
		return models.Restaurant{}, err
	}

	return restaurant, nil
}

func DeleteRestaurant(c *gin.Context, id string) (models.Restaurant, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Restaurant{}, err
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.Restaurant{}, errors.New("unauthorized")
	}

	var restaurant models.Restaurant
	err = restaurantCollection.FindOne(c.Request.Context(), bson.M{"_id": objectID}).Decode(&restaurant)
	if err != nil {
		return models.Restaurant{}, err
	}

	_, err = restaurantCollection.DeleteOne(c.Request.Context(), bson.M{"_id": objectID})
	if err != nil {
		return models.Restaurant{}, err
	}

	return restaurant, nil
}
