package controllers

import (
	"net/http"

	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/services"
	"github.com/gin-gonic/gin"
)

// @Summary Get all orders
// @Description Get all orders
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param 		 recordPerPage query int false "Record Per Page"
// @Param 		 page query int false "Page"
// @Param 		 startIndex query int false "Start Index"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /orders [get]
func GetOrders(c *gin.Context) {

	response, err := services.GetOrders(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": response.AllOrders, "page": response.Page, "recordPerPage": response.RecordPerPage, "startIndex": response.StartIndex})
}

// @Summary Get a order
// @Description Get a order
// @Tags User
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Order ID"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /order/{id} [get]
func GetOrder(c *gin.Context) {
	order_id := c.Param("id")

	order, err := services.GetOrderById(c, order_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order retrived successfully", "data": order, "status": http.StatusOK, "success": true})

}

// @Summary Create a order
// @Description Create a order
// @Tags User
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param table_id body string true "Table ID"
// @Param order_status body string true "Order Status"
// @Param total_amount body string true "Total Amount"
// @Success 201 {object} string
// @Failure 500 {object} string
// @Router /order [post]
func CreateOrder(c *gin.Context) {
	var orderReq models.Order

	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := services.CreateOrder(c, orderReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Order created successfully", "data": order, "status": http.StatusCreated, "success": true})
}

// @Summary Update a order
// @Description Update a order
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Order ID"
// @Param order body models.Order true "Table ID"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /order/{id} [put]
func UpdateOrder(c *gin.Context) {
	var reqOrder models.Order

	orderId := c.Param("id")
	if err := c.ShouldBindJSON(&reqOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := services.UpdateOrder(c, reqOrder, orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order updated successfully", "status": http.StatusOK, "success": true, "data": order})

}

// @Summary Delete a order
// @Description Delete a order
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Order ID"
// @Success 200 {object} string
// @Failure 500 {object} string
// @Router /order/{id} [delete]
func DeleteOrder(c *gin.Context) {
	orderId := c.Param("id")

	_, err := services.DeleteOrder(c, orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order deleted successfully", "status": http.StatusOK, "success": true, "data": nil})
}
