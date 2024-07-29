package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")

func GetUsers(c *gin.Context) {
	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, errp := strconv.Atoi(c.Query("page"))
	if errp != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(c.Query("startIndex"))
	if err != nil || startIndex < 0 {
		startIndex = 0
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}},
			},
		},
	}

	_, errAgg := userCollection.Aggregate(c.Request.Context(), mongo.Pipeline{matchStage, projectStage})

	if errAgg != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errAgg.Error()})
		return
	}

	var allUsers []bson.M

	c.JSON(http.StatusOK, gin.H{"data": allUsers[0], "total_count": allUsers[0]["total_count"], "page": page, "recordPerPage": recordPerPage, "startIndex": startIndex})

}

func GetUser(c *gin.Context) {
	userId := c.Param("id")

	if err := c.ShouldBindJSON(&userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	defer c.Request.Body.Close()

	err := userCollection.FindOne(c.Request.Context(), bson.D{{Key: "id", Value: userId}}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "status": http.StatusOK, "success": true, "error": false, "message": "User retrieved successfully"})

}

func UpdateUser(c *gin.Context) {
	fmt.Println("UpdateUser")
}

func DeleteUser(c *gin.Context) {
	userId := c.Param("id")

	if err := c.ShouldBindJSON(&userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := userCollection.DeleteOne(c.Request.Context(), bson.D{{Key: "id", Value: userId}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "User deleted successfully", "status": http.StatusOK, "success": true, "data": nil})
}

func Login(c *gin.Context) {
	var user models.User
	var foundUser models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := userCollection.FindOne(c.Request.Context(), bson.D{{Key: "email", Value: user.Email}}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	passwordIsValid, msg := ComparePassword(foundUser.Password, user.Password)

	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}

	token, refreshToken, _ := helpers.GenerateAllTokens(foundUser.User_id, foundUser.Email, foundUser.First_name, foundUser.Last_name)

	helpers.UpdateAllTokens(foundUser.User_id, token, refreshToken)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "User logged in successfully", "data": foundUser, "status": http.StatusOK, "success": true})

}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := userCollection.CountDocuments(c.Request.Context(), bson.D{{Key: "email", Value: user.Email}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	count, err = userCollection.CountDocuments(c.Request.Context(), bson.D{{Key: "phone", Value: user.Phone}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone number already exists"})
		return
	}

	password := HashPassword(user.Password)
	user.Password = password

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	token, refreshToken, _ := helpers.GenerateAllTokens(user.User_id, user.Email, user.First_name, user.Last_name)
	user.Token = token
	user.RefreshToken = refreshToken

	_, err = userCollection.InsertOne(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "User created successfully", "data": user, "status": http.StatusCreated, "success": true})

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
