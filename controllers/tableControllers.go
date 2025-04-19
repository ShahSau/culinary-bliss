package controllers

import (
	"fmt"
	"net/http"

	"github.com/ShahSau/culinary-bliss/services"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.GetCollection(database.DB, "tables")

// @Summary Get all tables
// @Description Get all tables
// @Tags Global
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /table [get]
func GetTables(c *gin.Context) {
	results, err := services.GetTables(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Table retrived successfully", "data": results, "status": http.StatusOK, "success": true})
}

// @Summary Get a table
// @Description Get a table
// @Tags Global
// @Accept json
// @Produce json
// @Param id path string true "Table ID"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /table/{id} [get]
func GetTable(c *gin.Context) {
	table_id := c.Param("id")

	table, err := services.GetTable(c, table_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Table retrived successfully", "data": table, "status": http.StatusOK, "success": true})
}

// @Summary Create a table
// @Description Create a table
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param table body types.Table true "Table"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Router /table [post]
func CreateTable(c *gin.Context) {
	var tableReq models.Table

	if err := c.ShouldBindJSON(&tableReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTable, err := services.CreateTable(c, tableReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Table created successfully", "data": newTable, "status": http.StatusCreated, "success": true})

}

// @Summary Update a table
// @Description Update a table
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Table ID"
// @Param table body types.Table true "Table"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /table/{id} [put]
func UpdateTable(c *gin.Context) {
	var tableReq models.Table
	id := c.Param("id")

	if err := c.ShouldBindJSON(&tableReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedTable, err := services.UpdateTable(c, id, tableReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Table updated successfully", "status": http.StatusOK, "success": true, "data": updatedTable})

}

// @Summary Delete a table
// @Description Delete a table
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Table ID"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /table/{id} [delete]
func DeleteTable(c *gin.Context) {
	id := c.Param("id")

	_, err := services.DeleteTable(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": fmt.Sprintf("Table with ID %s deleted successfully", id), "status": http.StatusOK, "success": true, "data": nil})
}
