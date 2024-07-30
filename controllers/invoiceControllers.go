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

type InvoiceViewFormat struct {
	Invoice_id        string      `json:"invoice_id"`
	Payment_method    string      `json:"payment_method"`
	Order_id          string      `json:"order_id"`
	Payment_status    string      `json:"payment_status"`
	Payment_due       interface{} `json:"payment_due"`
	Table_number      interface{} `json:"table_number"`
	Paymenet_due_date time.Time   `json:"payment_due_date"`
	Order_details     interface{} `json:"order_details"`
}

var invoiceCollection *mongo.Collection = database.GetCollection(database.DB, "invoice")

// @Summary Get Invoices
// @Description Get Invoices
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} models.Invoice
// @Failure 400 {object} string
// @Router /invoice [get]
func GetInvoices(c *gin.Context) {
	invoices, err := invoiceCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}

	defer invoices.Close(c.Request.Context())

	var results []InvoiceViewFormat

	for invoices.Next(c.Request.Context()) {
		var invoice InvoiceViewFormat
		if err = invoices.Decode(&invoice); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		}
		results = append(results, invoice)
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Invoice retrived successfully", "data": results, "status": http.StatusOK, "success": true})
}

// @Summary Get Invoice
// @Description Get Invoice
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} models.Invoice
// @Failure 400 {object} string
// @Router /invoice/{id} [get]
func GetInvoice(c *gin.Context) {
	var invoiceID = c.Param("id")

	var invoice models.Invoice

	defer c.Request.Body.Close()

	err := invoiceCollection.FindOne(c.Request.Context(), bson.M{"invoice_id": invoiceID}).Decode(&invoice)

	if err != nil {
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
// @Param invoice body types.Invoice true "Invoice"
// @Success 201 {object} models.Invoice
// @Failure 400 {object} string
// @Router /invoice [post]
func CreateInvoice(c *gin.Context) {
	var reqInvoice types.Invoice

	if err := c.ShouldBindJSON(&reqInvoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()
	var orderCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

	var order models.Order

	err := orderCollection.FindOne(c.Request.Context(), bson.M{"order_id": reqInvoice.Order_id}).Decode(&order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() + " Order not found"})
		return
	}

	var invoice models.Invoice
	invoice.Order_id = reqInvoice.Order_id
	invoice.Payment_method = reqInvoice.Payment_method
	invoice.Payment_status = "PENDING"
	invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.ID = primitive.NewObjectID()
	invoice.Invoice_id = invoice.ID.Hex()
	invoice.Total_amount = order.Total_amount

	_, err = invoiceCollection.InsertOne(c.Request.Context(), invoice)

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
// @Param id path string true "Invoice ID"
// @Param invoice body types.Invoice true "Invoice"
// @Success 200 {object} models.Invoice
// @Failure 400 {object} string
// @Router /invoice/{id} [put]
func UpdateInvoice(c *gin.Context) {
	var invoiceID = c.Param("id")

	var reqinvoice types.Invoice

	if err := c.ShouldBindJSON(&reqinvoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	var invoice models.Invoice
	err := invoiceCollection.FindOne(c.Request.Context(), bson.M{"invoice_id": invoiceID}).Decode(&invoice)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() + " Invoice not found"})
		return
	}

	var orderCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

	var order models.Order

	err = orderCollection.FindOne(c.Request.Context(), bson.M{"order_id": invoice.Order_id}).Decode(&order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() + " Order not found"})
		return
	}

	var updateObj models.Invoice

	updateObj.Order_id = reqinvoice.Order_id
	updateObj.Payment_method = reqinvoice.Payment_method
	updateObj.Payment_status = "PENDING"
	updateObj.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj.ID = invoice.ID
	updateObj.Invoice_id = invoice.Invoice_id
	updateObj.Total_amount = order.Total_amount

	_, err = invoiceCollection.UpdateOne(c.Request.Context(), bson.M{"invoice_id": invoiceID}, bson.D{{Key: "$set", Value: updateObj}})

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
// @Param id path string true "Invoice ID"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /invoice/{id} [delete]
func DeleteInvoice(c *gin.Context) {
	var invoiceID = c.Param("id")

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}

	_, err := invoiceCollection.DeleteOne(c.Request.Context(), bson.M{"invoice_id": invoiceID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Invoice deleted successfully", "status": http.StatusOK, "success": true, "data": nil})
}
