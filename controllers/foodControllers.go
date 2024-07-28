package controllers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.GetCollection(database.DB, "food")

// GetFoods godoc
func GetFoods(c *gin.Context) {

	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage

	_, err = strconv.Atoi(c.Query("startIndex"))

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
		}},
		{Key: "count", Value: bson.D{
			{Key: "$sum", Value: 1},
		}},
		{Key: "data", Value: bson.D{
			{Key: "$push", Value: "$$ROOT"},
		}},
	}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "count", Value: 1},
				{Key: "food_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			},
		},
	}

	_, errAgg := foodCollection.Aggregate(c.Request.Context(), mongo.Pipeline{matchStage, groupStage, projectStage})

	if errAgg != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var results []bson.M

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Food retrived successfully", "data": results[0], "status": http.StatusOK, "success": true})
}

func GetFood(c *gin.Context) {
	var foodID = c.Param("id")

	var food models.Food

	defer c.Request.Body.Close()

	err := foodCollection.FindOne(c.Request.Context(), primitive.M{"food_id": foodID}).Decode(&food)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"food":    food,
		"error":   false,
		"succes":  true,
		"message": "Food retrieved successfully",
		"status":  http.StatusOK,
	})
}

func CreateFood(c *gin.Context) {
	var reqfood models.Food

	var menu models.Menu
	var food models.Food

	if err := c.ShouldBindJSON(&reqfood); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// // Validate the input
	// if err := validate.Struct(reqfood); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	defer c.Request.Body.Close()
	reqfood.Food_id = primitive.NewObjectID().Hex()

	// Check if the menu exists
	err := database.GetCollection(database.DB, "menu").FindOne(c.Request.Context(), primitive.M{"menu_id": reqfood.Menu_id}).Decode(&menu)

	if err != nil {
		c.JSON(500, gin.H{"error": "Menu not found"})
		return
	}

	food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Food_id = reqfood.Food_id
	food.Name = reqfood.Name
	food.Description = reqfood.Description
	food.Price = toFixed(reqfood.Price, 2)
	food.Image = reqfood.Image
	food.Menu_id = reqfood.Menu_id
	food.ID = primitive.NewObjectID()

	foodCreated, err := foodCollection.InsertOne(c.Request.Context(), food)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"food":    foodCreated,
		"error":   false,
		"succes":  true,
		"message": "Food created successfully",
		"status":  http.StatusCreated,
	})
}

func UpdateFood(c *gin.Context) {
	var menu models.Menu
	var food models.Food

	foodId := c.Param("id")

	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateObj primitive.D

	if food.Name != "" {
		updateObj = append(updateObj, primitive.E{Key: "name", Value: food.Name})
	}

	if food.Description != "" {
		updateObj = append(updateObj, primitive.E{Key: "description", Value: food.Description})
	}

	if food.Price != 0 {
		updateObj = append(updateObj, primitive.E{Key: "price", Value: food.Price})
	}

	if food.Image != "" {
		updateObj = append(updateObj, primitive.E{Key: "image", Value: food.Image})
	}

	if food.Menu_id != "" {
		err := database.GetCollection(database.DB, "menu").FindOne(c.Request.Context(), primitive.M{"menu_id": food.Menu_id}).Decode(&menu)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu not found"})
			return
		}

		updateObj = append(updateObj, primitive.E{Key: "menu_id", Value: food.Menu_id})
	}

	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateObj = append(updateObj, primitive.E{Key: "updated_at", Value: food.UpdatedAt})

	_, err := foodCollection.UpdateOne(c.Request.Context(), bson.M{"food_id": foodId}, bson.D{{Key: "$set", Value: updateObj}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"succes":  true,
		"message": "Food updated successfully",
		"status":  http.StatusOK,
	})

}

func DeleteFood(c *gin.Context) {
	foodId := c.Param("id")

	_, err := foodCollection.DeleteOne(c.Request.Context(), bson.M{"food_id": foodId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"succes":  true,
		"message": "Food deleted successfully",
		"status":  http.StatusOK,
	})
}

func roundToTwo(num float64) float64 {
	return float64(int(num*100)) / 100
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output)) / output
}
