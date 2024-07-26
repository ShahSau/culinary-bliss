package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetOrderItems(c *gin.Context) {
	fmt.Println("GetOrderItems")
}

func GetOrderItem(c *gin.Context) {
	fmt.Println("GetOrderItem")
}

func CreateOrderItem(c *gin.Context) {
	fmt.Println("CreateOrderItem")
}

func UpdateOrderItem(c *gin.Context) {
	fmt.Println("UpdateOrderItem")
}

func DeleteOrderItem(c *gin.Context) {
	fmt.Println("DeleteOrderItem")
}
