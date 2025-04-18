package services

import (
	"errors"
	"log"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")

func LoginUser(user types.Loginuser, c *gin.Context) (models.User, string, string, error) {
	var foundUser models.User
	err := userCollection.FindOne(c.Request.Context(), bson.D{{Key: "email", Value: user.Email}}).Decode(&foundUser)
	if err != nil {
		return foundUser, "", "", errors.New("user not found")
	}

	passwordIsValid, msg := ComparePassword(foundUser.Password, user.Password)
	if !passwordIsValid {
		return foundUser, "", "", errors.New(msg)
	}

	token, refreshToken, _ := helpers.GenerateAllTokens(foundUser.User_id, foundUser.Email, foundUser.First_name, foundUser.Last_name)
	helpers.UpdateAllTokens(foundUser.User_id, token, refreshToken)

	return foundUser, token, refreshToken, nil
}

func RegisterUser(user models.User, c *gin.Context) (models.User, error) {
	count, err := userCollection.CountDocuments(nil, bson.D{{Key: "email", Value: user.Email}})
	if err != nil {
		return user, err
	}
	if count > 0 {
		return user, errors.New("email already exists")
	}

	count, err = userCollection.CountDocuments(c.Request.Context(), bson.D{{Key: "phone", Value: user.Phone}})
	if err != nil {
		return user, err
	}
	if count > 0 {
		return user, errors.New("phone number already exists")
	}

	password := HashPassword(user.Password)
	user.Password = password
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()
	user.Role = "User"

	token, refreshToken, _ := helpers.GenerateAllTokens(user.User_id, user.Email, user.First_name, user.Last_name)
	user.Token = token
	user.RefreshToken = refreshToken

	_, err = userCollection.InsertOne(nil, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func LogoutUser(userID string) error {
	helpers.UpdateAllTokens(userID, "", "")
	return nil
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic("Error hashing password")
	}
	return string(bytes)
}

func ComparePassword(hashedPassword string, password string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, "Invalid password"
	}
	return true, "Password is valid"
}
