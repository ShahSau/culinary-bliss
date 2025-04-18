package controllers

import (
	"math"
	"net/http"

	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/services"
	"github.com/gin-gonic/gin"
)

// @Summary GetFoods
// @Description Get all foods
// @Tags Global
// @Produce json
// @Param 		 recordPerPage query int false "Record Per Page"
// @Param 		 page query int false "Page"
// @Param 		 startIndex query int false "Start Index"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /foods [get]
func GetFoods(c *gin.Context) {
	response, err := services.GetFoods(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.AllFoods, "page": response.Page, "recordPerPage": response.RecordPerPage, "startIndex": response.StartIndex})

}

// @Summary Get Food
// @Description Get Food
// @Tags Global
// @Accept json
// @Produce json
// @Param id path string true "Food ID"
// @Success 200 {object}  string
// @Failure 400 {object} string
// @Router /food/{id} [get]
func GetFood(c *gin.Context) {
	foodId := c.Param("id")
	food, err := services.GetFoodByID(foodId, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data":    food,
		"error":   false,
		"succes":  true,
		"message": "Food retrieved successfully",
		"status":  http.StatusOK,
	})
}

// @Summary Create Food
// @Description Create Food
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param food body models.Food true "Food Object"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Router /food [post]
func CreateFood(c *gin.Context) {
	var reqfood models.Food
	if err := c.ShouldBindJSON(&reqfood); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	food, err := services.CreateFood(reqfood, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"data":    food,
		"error":   false,
		"succes":  true,
		"message": "Food created successfully",
		"status":  http.StatusCreated,
	})
}

// @Summary Update Food
// @Description Update Food
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Food ID"
// @Param food body models.Food true "Food Object"
// @Success 202 {object} string
// @Failure 400 {object} string
// @Router /food/{id} [put]
func UpdateFood(c *gin.Context) {
	var food models.Food

	foodId := c.Param("id")

	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateObj, err := services.UpdateFood(foodId, food, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"succes":  true,
		"message": "Food updated successfully",
		"status":  http.StatusOK,
		"data":    updateObj,
	})

}

// @Summary Delete Food
// @Description Delete Food
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Food ID"
// @Success 202 {object} string
// @Failure 400 {object} string
// @Router /food/{id} [delete]
func DeleteFood(c *gin.Context) {
	foodId := c.Param("id")

	_, err := services.DeleteFood(foodId, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"succes":  true,
		"message": "Food deleted successfully",
		"status":  http.StatusOK,
		"data":    nil,
	})
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output)) / output
}
