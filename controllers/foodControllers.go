package controllers

import (
	"fmt"

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
