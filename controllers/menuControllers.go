package controllers

import (
	"log"
	"net/http"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.GetCollection(database.DB, "menu")

// @Summary Get all menus
// @Description Get all menus
// @Tags Global
// @Accept json
// @Produce json
// @Param 		 recordPerPage query int false "Record Per Page"
// @Param 		 page query int false "Page"
// @Param 		 startIndex query int false "Start Index"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /menu [get]
func GetMenus(c *gin.Context) {
	response, err := services.GetMenus(c)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.AllMenus, "page": response.Page, "recordPerPage": response.RecordPerPage, "startIndex": response.StartIndex})
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

	menu, err := services.GetMenuByID(menuID, c)
	if err != nil {
		log.Fatal(err)
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
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param menu body types.Menu true "Menu object"
// @Success 201 {object} string
// @Failure 500 {object} string
// @Router /menu [post]
func CreateMenu(c *gin.Context) {
	var menu models.Menu

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqMenu, err := services.CreateMenu(menu, c)
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
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Menu ID"
// @Param menu body types.Menu true "Menu object"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /menu/{id} [put]
func UpdateMenu(c *gin.Context) {
	var reqMenu models.Menu

	if err := c.ShouldBindJSON(&reqMenu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menu, err := services.UpdateMenu(c.Param("id"), reqMenu, c)
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
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Menu ID"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /menu/{id} [delete]
func DeleteMenu(c *gin.Context) {
	menuId := c.Param("id")

	err := services.DeleteMenu(menuId, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Menu deleted successfully", "status": http.StatusOK, "success": true, "data": nil})
}
