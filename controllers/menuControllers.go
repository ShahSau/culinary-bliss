package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetMenus(c *gin.Context) {
	fmt.Println("GetMenus")
}

func GetMenu(c *gin.Context) {
	fmt.Println("GetMenu")
}

func CreateMenu(c *gin.Context) {
	fmt.Println("CreateMenu")
}

func UpdateMenu(c *gin.Context) {
	fmt.Println("UpdateMenu")
}

func DeleteMenu(c *gin.Context) {
	fmt.Println("DeleteMenu")
}
