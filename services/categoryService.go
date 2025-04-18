package services

import (
	"errors"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = database.GetCollection(database.DB, "categories")

func GetCategories(c *gin.Context) ([]models.Category, error) {
	var categories []models.Category
	cursor, err := categoryCollection.Find(c.Request.Context(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c.Request.Context())

	for cursor.Next(c.Request.Context()) {
		var category models.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetCategoryByID(id string, c *gin.Context) (models.Category, error) {
	var category models.Category
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return category, errors.New("invalid category ID")
	}

	err = categoryCollection.FindOne(c.Request.Context(), bson.M{"_id": objectID}).Decode(&category)
	if err != nil {
		return category, errors.New("category not found")
	}

	return category, nil
}

func CreateCategory(category models.Category, c *gin.Context) (models.Category, error) {

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return category, errors.New("unauthorized")
	}
	var newCategory models.Category
	newCategory.Title = category.Title
	newCategory.Image = category.Image
	newCategory.ID = primitive.NewObjectID()
	newCategory.Category_id = newCategory.ID.Hex()
	newCategory.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	newCategory.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err := categoryCollection.InsertOne(c.Request.Context(), newCategory)
	if err != nil {
		return category, err
	}

	return category, nil
}

func UpdateCategory(id string, updatedCategory models.Category, c *gin.Context) (models.Category, error) {
	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return updatedCategory, errors.New("invalid category ID")
	}
	var category models.Category

	category.Title = updatedCategory.Title
	category.Image = updatedCategory.Image
	category.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	category.CreatedAt = updatedCategory.CreatedAt
	category.ID = updatedCategory.ID
	category.Category_id = updatedCategory.Category_id
	update := bson.M{"$set": category}

	_, err = categoryCollection.UpdateOne(c.Request.Context(), bson.M{"_id": categoryID}, update)
	if err != nil {
		return category, err
	}

	return category, nil
}

func DeleteCategory(id string, c *gin.Context) error {
	categoryID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid category ID")
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return errors.New("unauthorized")
	}

	_, err = categoryCollection.DeleteOne(c.Request.Context(), bson.M{"_id": categoryID})
	if err != nil {
		return err
	}

	return nil
}
