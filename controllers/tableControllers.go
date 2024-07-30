package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/types"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	tables, err := tableCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer tables.Close(c.Request.Context())

	var results []models.Table

	for tables.Next(c.Request.Context()) {
		var table models.Table
		if err = tables.Decode(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		}
		results = append(results, table)
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

	var table models.Table

	defer c.Request.Body.Close()

	err := tableCollection.FindOne(c.Request.Context(), bson.M{"table_id": table_id}).Decode(&table)

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
	var tableReq types.Table

	if err := c.ShouldBindJSON(&tableReq); err != nil {
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

	var newTable models.Table

	newTable.Number_of_guests = tableReq.Number_of_guests
	newTable.Table_number = tableReq.Table_number
	newTable.Table_status = tableReq.Table_status
	newTable.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	newTable.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	newTable.ID = primitive.NewObjectID()
	newTable.Table_id = newTable.ID.Hex()

	_, err := tableCollection.InsertOne(c.Request.Context(), newTable)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Table created successfully", "data": tableReq, "status": http.StatusCreated, "success": true})

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
	var table types.Table
	table_id := c.Param("id")

	if err := c.ShouldBindJSON(&table); err != nil {
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

	var updatedTable models.Table

	err := tableCollection.FindOne(c.Request.Context(), bson.M{"table_id": table_id}).Decode(&updatedTable)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedTable.Number_of_guests = table.Number_of_guests
	updatedTable.Table_number = table.Table_number
	updatedTable.Table_status = table.Table_status
	updatedTable.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err = tableCollection.UpdateOne(c.Request.Context(), bson.M{"table_id": table_id}, bson.D{{Key: "$set", Value: updatedTable}})
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
	table_id := c.Param("id")

	_, err := tableCollection.DeleteOne(c.Request.Context(), bson.M{"table_id": table_id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": fmt.Sprintf("Table with ID %s deleted successfully", table_id), "status": http.StatusOK, "success": true, "data": nil})
}
