package controllers

import (
	"net/http"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.GetCollection(database.DB, "menu")

func GetMenus(c *gin.Context) {
	menus, err := menuCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer menus.Close(c.Request.Context())

	var results []models.Menu

	for menus.Next(c.Request.Context()) {
		var menu models.Menu
		if err = menus.Decode(&menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		}
		results = append(results, menu)
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Menu retrived successfully", "data": results, "status": http.StatusOK, "success": true})

}

func GetMenu(c *gin.Context) {
	var menuID = c.Param("id")

	var menu models.Menu

	defer c.Request.Body.Close()

	err := menuCollection.FindOne(c, primitive.M{"menu_id": menuID}).Decode(&menu)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"menu":    menu,
		"error":   false,
		"succes":  true,
		"message": "Menu retrieved successfully",
		"status":  http.StatusOK,
	})

}

func CreateMenu(c *gin.Context) {
	var reqMenu models.Menu

	if err := c.ShouldBindJSON(&reqMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the input
	if err := validate.Struct(reqMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqMenu.ID = primitive.NewObjectID()
	reqMenu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	defer c.Request.Body.Close()

	menuCreated, err := menuCollection.InsertOne(c.Request.Context(), reqMenu)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Menu created successfully", "data": menuCreated})
}

func UpdateMenu(c *gin.Context) {
	var reqMenu models.Menu

	if err := c.ShouldBindJSON(&reqMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the input
	if err := validate.Struct(reqMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menuId := c.Param("id")

	menu := models.Menu{}

	err := menuCollection.FindOne(c.Request.Context(), bson.M{"menu_id": menuId}).Decode(&menu)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	menu.ID = reqMenu.ID
	menu.Name = reqMenu.Name
	menu.Description = reqMenu.Description
	menu.Start_Date = reqMenu.Start_Date
	menu.End_Date = reqMenu.End_Date
	menu.Menu_id = menuId
	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err = menuCollection.UpdateOne(c.Request.Context(), bson.M{"menu_id": menuId}, menu)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Menu updated successfully", "data": menu, "status": http.StatusOK, "success": true})
}

func DeleteMenu(c *gin.Context) {
	menuId := c.Param("id")

	_, err := menuCollection.DeleteOne(c.Request.Context(), bson.M{"menu_id": menuId})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Menu deleted successfully"})
}
