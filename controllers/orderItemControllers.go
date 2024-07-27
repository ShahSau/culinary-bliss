package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetOrderItemsByOrder(c *gin.Context) {
	fmt.Println("GetOrderItemsByOrder")
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	return
}
