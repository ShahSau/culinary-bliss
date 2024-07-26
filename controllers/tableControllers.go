package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetTables(c *gin.Context) {
	fmt.Println("GetTables")
}

func GetTable(c *gin.Context) {
	fmt.Println("GetTable")
}

func CreateTable(c *gin.Context) {
	fmt.Println("CreateTable")
}

func UpdateTable(c *gin.Context) {
	fmt.Println("UpdateTable")
}

func DeleteTable(c *gin.Context) {
	fmt.Println("DeleteTable")
}
