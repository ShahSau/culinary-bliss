package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTables(c *gin.Context) ([]models.Table, error) {
	tables, err := tableCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		return nil, err
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

	return results, nil
}

func GetTable(c *gin.Context, id string) (models.Table, error) {

	table_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Table{}, err
	}

	var table models.Table

	defer c.Request.Body.Close()

	err = tableCollection.FindOne(c.Request.Context(), bson.M{"table_id": table_id}).Decode(&table)

	if err != nil {
		return table, err
	}

	return table, nil
}

func CreateTable(c *gin.Context, tableReq models.Table) (models.Table, error) {
	var table models.Table

	if err := c.ShouldBindJSON(&table); err != nil {
		return models.Table{}, err
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.Table{}, errors.New("Unauthorized")
	}

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
		return models.Table{}, err
	}

	return newTable, nil
}

func UpdateTable(c *gin.Context, id string, tableReq models.Table) (models.Table, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Table{}, err
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.Table{}, errors.New("unauthorized")
	}

	defer c.Request.Body.Close()

	var updatedTable models.Table

	err = tableCollection.FindOne(c.Request.Context(), bson.M{"table_id": objectID}).Decode(&updatedTable)

	if err != nil {
		return models.Table{}, err
	}

	updatedTable.Number_of_guests = tableReq.Number_of_guests
	updatedTable.Table_number = tableReq.Table_number
	updatedTable.Table_status = tableReq.Table_status
	updatedTable.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err = tableCollection.UpdateOne(c.Request.Context(), bson.M{"table_id": objectID}, bson.D{{Key: "$set", Value: updatedTable}})
	if err != nil {
		return models.Table{}, err
	}
	return updatedTable, nil
}

func DeleteTable(c *gin.Context, id string) (models.Table, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Table{}, err
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.Table{}, errors.New("unauthorized")
	}

	defer c.Request.Body.Close()

	var deletedTable models.Table

	err = tableCollection.FindOne(c.Request.Context(), bson.M{"table_id": objectID}).Decode(&deletedTable)

	if err != nil {
		return models.Table{}, err
	}

	_, err = tableCollection.DeleteOne(c.Request.Context(), bson.M{"table_id": objectID})
	if err != nil {
		return models.Table{}, err
	}
	return deletedTable, nil
}
