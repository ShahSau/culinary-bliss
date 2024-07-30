package controllers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")

// @Summary		Get all Users
// @Description	Get all users
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param 		 recordPerPage query int false "Record Per Page"
// @Param 		 page query int false "Page"
// @Param 		 startIndex query int false "Start Index"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/users [get]
func GetUsers(c *gin.Context) {
	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}
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

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "id", Value: 1},
				{Key: "first_name", Value: 1},
				{Key: "last_name", Value: 1},
				{Key: "email", Value: 1},
				{Key: "role", Value: 1},
				{Key: "phone", Value: 1},
			},
		},
	}

	record := bson.D{{Key: "$skip", Value: recordPerPage * (page - 1)}}
	limit := bson.D{{Key: "$limit", Value: recordPerPage}}

	result, errAgg := userCollection.Aggregate(c.Request.Context(), mongo.Pipeline{matchStage, projectStage, record, limit})

	if errAgg != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errAgg.Error()})
		return
	}

	var allUsers []bson.M
	if err = result.All(c.Request.Context(), &allUsers); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": allUsers, "page": page, "recordPerPage": recordPerPage, "startIndex": startIndex})

}

// @Summary		Get a  User
// @Description	Get a user Details by ID
// @Tags			User
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param 		 id path string true "User ID"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/users/{id} [get]
func GetUser(c *gin.Context) {
	userId := c.Param("id")

	var user models.User

	defer c.Request.Body.Close()

	err := userCollection.FindOne(c.Request.Context(), bson.D{{Key: "user_id", Value: userId}}).Decode(&user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "status": http.StatusOK, "success": true, "error": false, "message": "User retrieved successfully"})

}

// @Summary		Update User
// @Description	Update user details
// @Tags			User
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param 		 user body types.RegisterUser true "User"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/users/{id} [put]
func UpdateUser(c *gin.Context) {
	userId := c.Param("id")
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find the user
	var foundUser models.User
	err := userCollection.FindOne(c.Request.Context(), bson.D{{Key: "user_id", Value: userId}}).Decode(&foundUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = foundUser.Password
	user.CreatedAt = foundUser.CreatedAt
	user.ID = foundUser.ID
	user.Role = foundUser.Role
	user.User_id = foundUser.User_id
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err = userCollection.UpdateOne(c.Request.Context(), bson.D{{Key: "user_id", Value: userId}}, bson.D{{Key: "$set", Value: user}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "User updated successfully", "status": http.StatusOK, "success": true, "data": user})
}

// @Summary		Delete User
// @Description	Delete user
// @Tags			Admin
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param 		 id path string true "User ID"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	userId := c.Param("id")

	defer c.Request.Body.Close()

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}

	_, err := userCollection.DeleteOne(c.Request.Context(), bson.D{{Key: "id", Value: userId}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "User deleted successfully", "status": http.StatusOK, "success": true, "data": nil})
}

// @Summary		Reset Password
// @Description	Reset user password
// @Tags			User
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param 		 user body types.PasswordReset true "User"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/reset-password [post]
func ResetPassword(c *gin.Context) {
	var user types.PasswordReset
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var foundUser models.User
	err := userCollection.FindOne(c.Request.Context(), bson.D{{Key: "email", Value: user.Email}}).Decode(&foundUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	passwordIsValid, msg := ComparePassword(foundUser.Password, user.OldPassword)

	if !passwordIsValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
		return
	}

	new_password := HashPassword(user.NewPassword)
	foundUser.Password = new_password

	foundUser.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err = userCollection.UpdateOne(c.Request.Context(), bson.D{{Key: "email", Value: user.Email}}, bson.D{{Key: "$set", Value: foundUser}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Password reset successfully", "status": http.StatusOK, "success": true, "data": foundUser})
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
