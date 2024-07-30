package controllers

import (
	"log"
	"math"
	"net/http"
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

// @Summary GetFoods
// @Description Get all foods
// @Tags Global
// @Produce json
// @Param 		 recordPerPage query int false "Record Per Page"
// @Param 		 page query int false "Page"
// @Param 		 startIndex query int false "Start Index"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /foods [get]
func GetFoods(c *gin.Context) {
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
				{Key: "name", Value: 1},
				{Key: "image", Value: 1},
				{Key: "description", Value: 1},
				{Key: "price", Value: 1},
				{Key: "menu_id", Value: 1},
			},
		},
	}

	record := bson.D{{Key: "$skip", Value: recordPerPage * (page - 1)}}
	limit := bson.D{{Key: "$limit", Value: recordPerPage}}

	result, errAgg := foodCollection.Aggregate(c.Request.Context(), mongo.Pipeline{matchStage, projectStage, record, limit})

	if errAgg != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errAgg.Error()})
		return
	}

	var allFoods []bson.M
	if err = result.All(c.Request.Context(), &allFoods); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": allFoods, "page": page, "recordPerPage": recordPerPage, "startIndex": startIndex})

}

// @Summary Get Food
// @Description Get Food
// @Tags Global
// @Accept json
// @Produce json
// @Param id path string true "Food ID"
// @Success 200 {object}  string
// @Failure 400 {object} string
// @Router /food/{id} [get]
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
		"data":    food,
		"error":   false,
		"succes":  true,
		"message": "Food retrieved successfully",
		"status":  http.StatusOK,
	})
}

// @Summary Create Food
// @Description Create Food
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param food body models.Food true "Food Object"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Router /food [post]
func CreateFood(c *gin.Context) {
	var reqfood models.Food

	var menu models.Menu
	var food models.Food

	if err := c.ShouldBindJSON(&reqfood); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}

	defer c.Request.Body.Close()
	reqfood.Food_id = primitive.NewObjectID().Hex()

	// Check if the menu exists
	err := database.GetCollection(database.DB, "menu").FindOne(c.Request.Context(), primitive.M{"menu_id": reqfood.Menu_id}).Decode(&menu)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu not found"})
		return
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data":    food,
		"error":   false,
		"succes":  true,
		"message": "Food created successfully",
		"status":  http.StatusCreated,
	})
}

// @Summary Update Food
// @Description Update Food
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Food ID"
// @Param food body models.Food true "Food Object"
// @Success 202 {object} string
// @Failure 400 {object} string
// @Router /food/{id} [put]
func UpdateFood(c *gin.Context) {
	var menu models.Menu
	var food models.Food

	foodId := c.Param("id")

	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
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
		"data":    updateObj,
	})

}

// @Summary Delete Food
// @Description Delete Food
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Food ID"
// @Success 202 {object} string
// @Failure 400 {object} string
// @Router /food/{id} [delete]
func DeleteFood(c *gin.Context) {
	foodId := c.Param("id")

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}
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
		"data":    nil,
	})
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output)) / output
}
