package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {
	fmt.Println("GetOrders")
}

func GetOrder(c *gin.Context) {
	fmt.Println("GetOrder")
}

func CreateOrder(c *gin.Context) {
	fmt.Println("CreateOrder")
}

func UpdateOrder(c *gin.Context) {
	fmt.Println("UpdateOrder")
}

func DeleteOrder(c *gin.Context) {
	fmt.Println("DeleteOrder")
}
