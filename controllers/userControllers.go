package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	fmt.Println("GetUsers")
}

func GetUser(c *gin.Context) {
	fmt.Println("GetUser")
}

func CreateUser(c *gin.Context) {
	fmt.Println("CreateUser")
}

func UpdateUser(c *gin.Context) {
	fmt.Println("UpdateUser")
}

func DeleteUser(c *gin.Context) {
	fmt.Println("DeleteUser")
}
