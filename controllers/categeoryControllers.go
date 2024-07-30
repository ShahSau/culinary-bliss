package controllers

import (
	"net/http"
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

var categoryCollection *mongo.Collection = database.GetCollection(database.DB, "categories")

// @Summary Get all categories
// @Description Get all categories
// @Tags Global
// @Accept json
// @Produce json
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/categories [get]
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

// @Summary Get a category
// @Description Get a category
// @Tags User
// @Accept json
// @Produce json
// @param id path string true "Category ID"
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/categeory/{id} [get]
func GetCategory(c *gin.Context) {
	category_id := c.Param("id")

	var category models.Category

	defer c.Request.Body.Close()

	err := categoryCollection.FindOne(c.Request.Context(), bson.M{"category_id": category_id}).Decode(&category)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Category retrived successfully", "data": category, "status": http.StatusOK, "success": true})
}

// @Summary Create a category
// @Description Create a category
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @param Authorization header string true "Token"
// @Param category body types.Category true "Category"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/categories [post]
func CreateCategory(c *gin.Context) {
	var category types.Category

	if err := c.ShouldBindJSON(&category); err != nil {
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

	var newCategory models.Category

	newCategory.Title = category.Title
	newCategory.Image = category.Image
	newCategory.ID = primitive.NewObjectID()
	newCategory.Category_id = newCategory.ID.Hex()
	newCategory.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	newCategory.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err := categoryCollection.InsertOne(c.Request.Context(), newCategory)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Category created successfully", "status": http.StatusOK, "success": true, "data": category})
}

// @Summary Update a category
// @Description Update a category
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @param Authorization header string true "Token"
// @param id path string true "Category ID"
// @Param category body types.Category true "Category"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/categeory/{id} [put]
func UpdateCategory(c *gin.Context) {
	category_id := c.Param("id")
	var req_category types.Category
	//var category models.Category

	if err := c.ShouldBindJSON(&req_category); err != nil {
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

	// find the category
	var category models.Category
	err := categoryCollection.FindOne(c.Request.Context(), bson.M{"category_id": category_id}).Decode(&category)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updatedCategory models.Category

	updatedCategory.Title = req_category.Title
	updatedCategory.Image = req_category.Image
	updatedCategory.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updatedCategory.CreatedAt = category.CreatedAt
	updatedCategory.ID = category.ID
	updatedCategory.Category_id = category.Category_id

	_, err = categoryCollection.UpdateOne(c.Request.Context(), bson.M{"category_id": category_id}, bson.D{{Key: "$set", Value: updatedCategory}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Category updated successfully", "status": http.StatusOK, "success": true, "data": category})

}

// @Summary Delete a category
// @Description Delete a category
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @param Authorization header string true "Token"
// @param id path string true "Category ID"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/categeory/{id} [delete]
func DeleteCategory(c *gin.Context) {
	category_id := c.Param("id")

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
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
