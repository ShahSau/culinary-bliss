package services

import (
	"errors"
	"math"
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

var foodCollection *mongo.Collection = database.GetCollection(database.DB, "food")

func GetFoods(c *gin.Context) (models.Response, error) {
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
		{Key: "price", Value: 1},
		{Key: "menu_id", Value: 1},
	}}}
	skipStage := bson.D{{Key: "$skip", Value: startIndex}}
	limitStage := bson.D{{Key: "$limit", Value: recordPerPage}}

	cursor, err := foodCollection.Aggregate(c.Request.Context(), mongo.Pipeline{
		matchStage, projectStage, skipStage, limitStage,
	})
	if err != nil {
		return models.Response{}, err
	}
	defer cursor.Close(c.Request.Context())

	var foods []models.Food
	for cursor.Next(c.Request.Context()) {
		var food models.Food
		if err := cursor.Decode(&food); err != nil {
			return models.Response{}, err
		}
		foods = append(foods, food)
	}

	response := models.Response{
		AllFoods:      foods,
		Page:          page,
		RecordPerPage: recordPerPage,
		StartIndex:    startIndex,
	}

	return response, nil
}

func GetFoodByID(id string, c *gin.Context) (models.Food, error) {
	var food models.Food
	foodID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return food, errors.New("invalid food ID")
	}

	err = foodCollection.FindOne(c.Request.Context(), bson.M{"_id": foodID}).Decode(&food)
	if err != nil {
		return food, errors.New("food not found")
	}

	return food, nil
}

func CreateFood(reqfood models.Food, c *gin.Context) (models.Food, error) {
	var menu models.Menu
	var food models.Food
	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return reqfood, errors.New("unauthorized")
	}

	// Check if the menu exists
	err := database.GetCollection(database.DB, "menu").FindOne(c.Request.Context(), primitive.M{"menu_id": reqfood.Menu_id}).Decode(&menu)

	if err != nil {
		return reqfood, errors.New("menu not found")
	}

	food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Name = reqfood.Name
	food.Description = reqfood.Description
	food.Price = toFixed(reqfood.Price, 2)
	food.Image = reqfood.Image
	food.Menu_id = reqfood.Menu_id
	food.ID = primitive.NewObjectID()
	food.Food_id = food.ID.Hex()

	_, err = foodCollection.InsertOne(c.Request.Context(), food)
	if err != nil {
		return food, err
	}

	return food, nil
}

func UpdateFood(id string, reqfood models.Food, c *gin.Context) (models.Food, error) {
	var menu models.Menu
	foodID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return reqfood, errors.New("invalid food ID")
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return reqfood, errors.New("unauthorized")
	}

	var updateObj primitive.D

	if reqfood.Name != "" {
		updateObj = append(updateObj, primitive.E{Key: "name", Value: reqfood.Name})
	}

	if reqfood.Description != "" {
		updateObj = append(updateObj, primitive.E{Key: "description", Value: reqfood.Description})
	}

	if reqfood.Price != 0 {
		updateObj = append(updateObj, primitive.E{Key: "price", Value: reqfood.Price})
	}

	if reqfood.Image != "" {
		updateObj = append(updateObj, primitive.E{Key: "image", Value: reqfood.Image})
	}

	if reqfood.Menu_id != "" {
		err := database.GetCollection(database.DB, "menu").FindOne(c.Request.Context(), primitive.M{"menu_id": reqfood.Menu_id}).Decode(&menu)

		if err != nil {
			return reqfood, errors.New("menu not found")
		}

		updateObj = append(updateObj, primitive.E{Key: "menu_id", Value: reqfood.Menu_id})
	}

	reqfood.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateObj = append(updateObj, primitive.E{Key: "updated_at", Value: reqfood.UpdatedAt})

	_, err = foodCollection.UpdateOne(c.Request.Context(), bson.M{"food_id": foodID}, bson.D{{Key: "$set", Value: updateObj}})

	if err != nil {
		return reqfood, err
	}

	return reqfood, nil
}

func DeleteFood(id string, c *gin.Context) (models.Food, error) {
	var food models.Food
	foodID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return food, errors.New("invalid food ID")
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return food, errors.New("unauthorized")
	}

	result, err := foodCollection.DeleteOne(c.Request.Context(), bson.M{"_id": foodID})
	if err != nil {
		return food, err
	}

	if result.DeletedCount == 0 {
		return food, errors.New("food not found")
	}

	return food, nil
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output)) / output
}
