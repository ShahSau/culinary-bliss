package controllers

import (
	"net/http"

	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/services"
	"github.com/ShahSau/culinary-bliss/types"
	"github.com/gin-gonic/gin"
)

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
	response, err := services.GetUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": response.AllUsers, "page": response.Page, "recordPerPage": response.RecordPerPage, "startIndex": response.StartIndex, "status": http.StatusOK, "success": true, "error": false, "message": "Users retrieved successfully"})

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

	user, err := services.GetUser(c, userId)
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
	var userReq models.User
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := services.UpdateUser(c, userId, userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "User updated successfully", "status": http.StatusOK, "success": true, "data": updatedUser})
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
	id := c.Param("id")

	_, err := services.DeleteUser(c, id)
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
	var userReq types.PasswordReset
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundUser, err := services.ResetPassword(c, userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Password reset successfully", "status": http.StatusOK, "success": true, "data": foundUser})
}
