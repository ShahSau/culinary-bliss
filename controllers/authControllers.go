package controllers

import (
	"net/http"
	"time"

	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary		User Login
// @Description	user can login by giving their email and password
// @Tags			Auth
// @Accept			json
// @Produce		    json
// @Param           user body types.Loginuser true "User"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/login [post]
func Login(c *gin.Context) {
	var user types.Loginuser
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

// @Summary		User Signup
// @Description	user can signup by giving their details
// @Tags			Auth
// @Accept			json
// @Produce		    json
// @Param 		 user body types.RegisterUser true "User"
// @Success		201	{object}	string
// @Failure		500	{object}	string
// @Router			/register [post]
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
	user.Role = "User"

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

// @Summary		User Logout
// @Description	user can logout by giving their user_id
// @Tags			Auth
// @Accept			json
// @Produce		    json
// @Param			user_id body string true "User ID"
// @Success		200	{object}	string
// @Failure		500	{object}	string
// @Router			/logout [post]
func Logout(c *gin.Context) {
	var user struct {
		User_id string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	helpers.UpdateAllTokens(user.User_id, "", "")

	c.Set("email", "")
	c.Set("first_name", "")
	c.Set("last_name", "")
	c.Set("user_id", "")

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "User logged out successfully", "status": http.StatusOK, "success": true})
}
