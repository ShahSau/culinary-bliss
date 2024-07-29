package controllers

import (
	"net/http"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = database.GetCollection(database.DB, "categories")

func GetCategories(c *gin.Context) {
	categories, err := categoryCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer categories.Close(c.Request.Context())

	var results []models.Category

	for categories.Next(c.Request.Context()) {
		var category models.Category
		if err = categories.Decode(&category); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		results = append(results, category)
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Categories retrived successfully", "data": results, "status": http.StatusOK, "success": true})
}

func GetCategory(c *gin.Context) {
	category_id := c.Param("id")

	if err := c.ShouldBindJSON(&category_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category

	defer c.Request.Body.Close()

	err := categoryCollection.FindOne(c.Request.Context(), bson.M{"category_id": category_id}).Decode(&category)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Category retrived successfully", "data": category, "status": http.StatusOK, "success": true})
}

func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	_, err := categoryCollection.InsertOne(c.Request.Context(), category)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Category created successfully", "status": http.StatusOK, "success": true, "data": category})
}

func UpdateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	category_id := c.Param("id")

	_, err := categoryCollection.UpdateOne(c.Request.Context(), bson.M{"category_id": category_id}, bson.D{{Key: "$set", Value: category}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Category updated successfully", "status": http.StatusOK, "success": true, "data": category})

}

func DeleteCategory(c *gin.Context) {
	category_id := c.Param("id")

	if err := c.ShouldBindJSON(&category_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()
	_, err := categoryCollection.DeleteOne(c.Request.Context(), bson.M{"category_id": category_id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Category deleted successfully", "status": http.StatusOK, "success": true})
}
