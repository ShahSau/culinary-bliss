package controllers

import (
	"fmt"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.GetCollection(database.DB, "tables")

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
