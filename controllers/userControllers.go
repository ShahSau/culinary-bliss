package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	c.JSON(http.StatusOK, gin.H{"data": user, "status": http.StatusOK, "success": true})

}

func CreateUser(c *gin.Context) {
	fmt.Println("CreateUser")
}

func UpdateUser(c *gin.Context) {
	fmt.Println("UpdateUser")
}

func DeleteUser(c *gin.Context) {
	fmt.Println("DeleteUser")
}

func Login(c *gin.Context) {
	// convert json to struct
	// find the user using email
	// compare password
	// generate token and referesh token
	// update the token in the database

}

func Register(c *gin.Context) {
	// convert json to struct
	// email exists?
	// phone number exists?
	// hash password

	//generate token and referesh token
	// insert into database
}

func HashPassword(password string) string {
	return password
}

func ComparePassword(hashedPassword string, password string) bool {
	fmt.Println("ComparePassword")
	return true
}
