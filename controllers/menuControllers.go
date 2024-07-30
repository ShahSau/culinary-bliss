package controllers

import (
	"net/http"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.GetCollection(database.DB, "menu")

// @Summary Get all menus
// @Description Get all menus
// @Tags Global
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /menu [get]
func GetMenus(c *gin.Context) {
	menus, err := menuCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Menu retrived successfully", "data": results, "status": http.StatusOK, "success": true})

}

// @Summary Get a menu
// @Description Get a menu
// @Tags Global
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /menu/{id} [get]
func GetMenu(c *gin.Context) {
	var menuID = c.Param("id")

	var menu models.Menu

	defer c.Request.Body.Close()

	err := menuCollection.FindOne(c.Request.Context(), primitive.M{"menu_id": menuID}).Decode(&menu)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    menu,
		"error":   false,
		"succes":  true,
		"message": "Menu retrieved successfully",
		"status":  http.StatusOK,
	})

}

// @Summary Create a menu
// @Description Create a menu
// @Tags Admin
// @Accept json
// @Produce json
// @Param menu body types.Menu true "Menu object"
// @Success 201 {object} string
// @Failure 500 {object} string
// @Router /menu [post]
func CreateMenu(c *gin.Context) {
	var incomingMenu types.Menu

	if err := c.ShouldBindJSON(&incomingMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reqMenu models.Menu

	reqMenu.Name = incomingMenu.Name
	reqMenu.Description = incomingMenu.Description
	reqMenu.Start_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.End_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.ID = primitive.NewObjectID()
	reqMenu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.Menu_id = reqMenu.ID.Hex()

	defer c.Request.Body.Close()

	_, err := menuCollection.InsertOne(c.Request.Context(), reqMenu)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Menu created successfully", "data": reqMenu, "status": http.StatusCreated, "success": true})
}

// @Summary Update a menu
// @Description Update a menu
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Param menu body types.Menu true "Menu object"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /menu/{id} [put]
func UpdateMenu(c *gin.Context) {
	var reqMenu types.Menu

	if err := c.ShouldBindJSON(&reqMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}

	menuId := c.Param("id")

	var menu models.Menu

	err := menuCollection.FindOne(c.Request.Context(), bson.M{"menu_id": menuId}).Decode(&menu)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	menu.Name = reqMenu.Name
	menu.Description = reqMenu.Description
	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.Start_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.End_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err = menuCollection.UpdateOne(c.Request.Context(), bson.M{"menu_id": menuId}, bson.M{"$set": menu})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Menu updated successfully", "data": menu, "status": http.StatusOK, "success": true})
}

// @Summary Delete a menu
// @Description Delete a menu
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /menu/{id} [delete]
func DeleteMenu(c *gin.Context) {
	menuId := c.Param("id")

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}

	_, err := menuCollection.DeleteOne(c.Request.Context(), bson.M{"menu_id": menuId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Menu deleted successfully", "status": http.StatusOK, "success": true, "data": nil})
}
