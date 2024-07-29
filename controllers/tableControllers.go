package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.GetCollection(database.DB, "tables")

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

func GetTable(c *gin.Context) {
	table_id := c.Param("id")

	if err := c.ShouldBindJSON(&table_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var table models.Table

	defer c.Request.Body.Close()

	err := tableCollection.FindOne(c.Request.Context(), bson.M{"table_id": table_id}).Decode(&table)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Table retrived successfully", "data": table, "status": http.StatusOK, "success": true})
}

func CreateTable(c *gin.Context) {
	var tableReq models.Table

	if err := c.ShouldBindJSON(&tableReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	tableReq.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	tableReq.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	tableReq.ID = primitive.NewObjectID()
	tableReq.Table_id = tableReq.ID.Hex()

	_, err := tableCollection.InsertOne(c.Request.Context(), tableReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Table created successfully", "data": tableReq, "status": http.StatusCreated, "success": true})

}

func UpdateTable(c *gin.Context) {
	var table models.Table
	table_id := c.Param("id")

	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateObj primitive.D

	if table.Number_of_guests != 0 {
		updateObj = append(updateObj, bson.E{Key: "number_of_guests", Value: table.Number_of_guests})
	}

	if table.Table_number != 0 {
		updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.Table_number})
	}

	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.UpdatedAt})

	_, err := tableCollection.UpdateOne(c.Request.Context(), bson.M{"table_id": table_id}, bson.D{{Key: "$set", Value: updateObj}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Table updated successfully", "status": http.StatusOK, "success": true})

}

func DeleteTable(c *gin.Context) {
	table_id := c.Param("id")

	if err := c.ShouldBindJSON(&table_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := tableCollection.DeleteOne(c.Request.Context(), bson.M{"table_id": table_id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": fmt.Sprintf("Table with ID %s deleted successfully", table_id), "status": http.StatusOK, "success": true, "data": nil})
}
