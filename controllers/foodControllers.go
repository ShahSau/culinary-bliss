package controllers

import (
	"fmt"
	"math"

	"github.com/gin-gonic/gin"
)

// GetFoods godoc
func GetFoods(c *gin.Context) {
	fmt.Println("GetFoods")
}

func GetFood(c *gin.Context) {
	fmt.Println("GetFood")
}

func CreateFood(c *gin.Context) {
	fmt.Println("CreateFood")
}

func UpdateFood(c *gin.Context) {
	fmt.Println("UpdateFood")
}

func DeleteFood(c *gin.Context) {
	fmt.Println("DeleteFood")
}

func roundToTwo(num float64) float64 {
	return float64(int(num*100)) / 100
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(int(num*output)) / output
}
