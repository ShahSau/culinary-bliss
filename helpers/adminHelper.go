package helpers

import (
	"context"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsAdmin(email string) bool {
	var user *mongo.Collection = database.GetCollection(database.DB, "users")
	var result models.User
	err := user.FindOne(context.Background(), bson.M{"email": email}).Decode(&result)
	if err != nil {
		return false
	}
	if result.Role == "Admin" {
		return true
	}
	return false
}
