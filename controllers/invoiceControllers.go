package controllers

import (
	"net/http"

	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/services"
	"github.com/gin-gonic/gin"
)

// @Summary Get Invoices
// @Description Get Invoices
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Success 200 {object} models.Invoice
// @Failure 400 {object} string
// @Router /invoice [get]
func GetInvoices(c *gin.Context) {
	results, err := services.GetInvoices(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Invoice retrived successfully", "data": results, "status": http.StatusOK, "success": true})
}

// @Summary Get Invoice
// @Description Get Invoice
// @Tags User
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Invoice ID"
// @Success 200 {object} models.Invoice
// @Failure 400 {object} string
// @Router /invoice/{id} [get]
func GetInvoice(c *gin.Context) {
	var invoiceID = c.Param("id")

	invoice, err := services.GetInvoiceByID(c, invoiceID)
	if err != nil {
		if err.Error() == "invoice not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Invoice retrieved successfully", "data": invoice, "status": http.StatusOK, "success": true})
}

// @Summary Create Invoice
// @Description Create Invoice
// @Tags User
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param invoice body types.Invoice true "Invoice"
// @Success 201 {object} models.Invoice
// @Failure 400 {object} string
// @Router /invoice [post]
func CreateInvoice(c *gin.Context) {
	var reqInvoice models.Invoice
	if err := c.ShouldBindJSON(&reqInvoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoice, err := services.CreateInvoice(c, reqInvoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Invoice created successfully", "status": http.StatusCreated, "success": true, "data": invoice})

}

// @Summary Update Invoice
// @Description Update Invoice
// @Tags User
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Invoice ID"
// @Param invoice body types.Invoice true "Invoice"
// @Success 200 {object} models.Invoice
// @Failure 400 {object} string
// @Router /invoice/{id} [put]
func UpdateInvoice(c *gin.Context) {
	var invoiceID = c.Param("id")
	var reqinvoice models.Invoice

	if err := c.ShouldBindJSON(&reqinvoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateObj, err := services.UpdateInvoice(c, invoiceID, reqinvoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Invoice updated successfully", "status": http.StatusOK, "success": true, "data": updateObj})

}

// @Summary Delete Invoice
// @Description Delete Invoice
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Invoice ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /invoice/{id} [delete]
func DeleteInvoice(c *gin.Context) {
	var invoiceID = c.Param("id")

	err := services.DeleteInvoice(c, invoiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Invoice deleted successfully", "status": http.StatusOK, "success": true, "data": nil})
}
