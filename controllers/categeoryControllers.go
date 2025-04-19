package controllers

import (
	"net/http"

	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/services"
	"github.com/gin-gonic/gin"
)

// @Summary Get all categories
// @Description Get all categories
// @Tags Global
// @Accept json
// @Produce json
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/categories [get]
func GetCategories(c *gin.Context) {
	categories, err := services.GetCategories(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// @Summary Get a category
// @Description Get a category
// @Tags User
// @Accept json
// @Produce json
// @param id path string true "Category ID"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/categeory/{id} [get]
func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	category, err := services.GetCategoryByID(id, c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
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
	id := c.Param("id")
	err := services.DeleteCategory(id, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
