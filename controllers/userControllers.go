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

func Login(c *gin.Context) {
	fmt.Println("Login")
}

func Register(c *gin.Context) {
	fmt.Println("Register")
}

func HashPassword(password string) string {
	return password
}

func ComparePassword(hashedPassword string, password string) bool {
	fmt.Println("ComparePassword")
	return true
}
