package controllers

import (
	"net/http"

	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/services"
	"github.com/ShahSau/culinary-bliss/types"
	"github.com/gin-gonic/gin"
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
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundUser, token, refreshToken, err := services.LoginUser(user, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "User logged in successfully", "data": foundUser, "token": token, "refreshToken": refreshToken, "status": http.StatusOK, "success": true})
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

	createdUser, err := services.RegisterUser(user, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "User created successfully", "data": createdUser, "status": http.StatusCreated, "success": true})
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

	err := services.LogoutUser(user.User_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "User logged out successfully", "status": http.StatusOK, "success": true})
}
