package controllers

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/bluesuncorp/validator.v5"
)

var foodCollection *mongo.Collection = database.GetCollection(database.DB, "food")
var validate = validator.New("validate", validator.BakedInValidators)

// GetFoods godoc
func GetFoods(c *gin.Context) {
	fmt.Println("GetFoods")
}

func GetFood(c *gin.Context) {
	var foodID = c.Param("id")

	var food models.Food

	defer c.Request.Body.Close()

	err := foodCollection.FindOne(c, primitive.M{"food_id": foodID}).Decode(&food)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
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

	// Validate the input
	if err := validate.Struct(reqfood); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()
	reqfood.Food_id = primitive.NewObjectID().Hex()

	// Check if the menu exists
	err := database.GetCollection(database.DB, "menu").FindOne(c, primitive.M{"menu_id": reqfood.Menu_id}).Decode(&menu)

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

	foodCreated, err := foodCollection.InsertOne(c, food)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"food":    foodCreated,
		"error":   false,
		"succes":  true,
		"message": "Food created successfully",
		"status":  http.StatusCreated,
	})
}

func UpdateFood(c *gin.Context) {
	fmt.Println("UpdateFood")
}

func DeleteFood(c *gin.Context) {
	fmt.Println("DeleteFood")
}

func roundToTwo(num float64) float64 {
	return float64(int(num*100)) / 100
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output)) / output
}
