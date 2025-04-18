package controllers

import (
	"net/http"

	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/services"
	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	categories, err := services.GetCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := services.GetCategoryByID(id, c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdCategory, err := services.CreateCategory(category, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully", "data": createdCategory})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := services.UpdateCategory(id, updatedCategory, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully", "data": category})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	err := services.DeleteCategory(id, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
