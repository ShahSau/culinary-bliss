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

func GetInvoices(c *gin.Context) {
	invoices, err := invoiceCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Invoice retrived successfully", "data": results, "status": http.StatusOK, "success": true})
}

func GetInvoice(c *gin.Context) {
	var invoiceID = c.Param("id")

	var invoice models.Invoice

	defer c.Request.Body.Close()

	err := invoiceCollection.FindOne(c.Request.Context(), bson.M{"invoice_id": invoiceID}).Decode(&invoice)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var invoiceView InvoiceViewFormat

	allOrders, err := ItemsByOrder(invoice.Order_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	invoiceView.Order_id = invoice.Order_id
	invoiceView.Paymenet_due_date = invoice.Payment_due_date

	invoiceView.Payment_method = invoice.Payment_method
	invoiceView.Payment_status = invoice.Payment_status
	invoiceView.Invoice_id = invoice.Invoice_id
	invoiceView.Payment_due = allOrders[0]["payment_due"]
	invoiceView.Table_number = allOrders[0]["table_number"]
	invoiceView.Order_details = allOrders[0]["order_items"]

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Invoice retrieved successfully", "data": invoiceView, "status": http.StatusOK, "success": true})
}

func CreateInvoice(c *gin.Context) {
	var invoice models.Invoice

	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()
	var orderCollection *mongo.Collection = database.GetCollection(database.DB, "order")

	var order models.Order

	err := orderCollection.FindOne(c.Request.Context(), bson.M{"order_id": invoice.Order_id}).Decode(&order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	invoice.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoice.ID = primitive.NewObjectID()
	invoice.Invoice_id = invoice.ID.Hex()

	_, err = invoiceCollection.InsertOne(c.Request.Context(), invoice)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": nil, "message": "Invoice created successfully", "status": http.StatusCreated, "success": true})

}

func UpdateInvoice(c *gin.Context) {
	var invoiceID = c.Param("id")

	var invoice models.Invoice

	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	err := invoiceCollection.FindOne(c.Request.Context(), bson.M{"invoice_id": invoiceID}).Decode(&invoice)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updateObj primitive.D

	if invoice.Payment_method != "" {
		updateObj = append(updateObj, bson.E{Key: "payment_method", Value: invoice.Payment_method})
	}

	if invoice.Payment_status != "" {
		updateObj = append(updateObj, bson.E{Key: "payment_status", Value: invoice.Payment_status})
	}

	if invoice.Order_id != "" {
		updateObj = append(updateObj, bson.E{Key: "order_id", Value: invoice.Order_id})
	}

	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})

	_, err = invoiceCollection.UpdateOne(c.Request.Context(), bson.M{"invoice_id": invoiceID}, bson.D{{Key: "$set", Value: updateObj}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Invoice updated successfully", "status": http.StatusOK, "success": true})

}

func DeleteInvoice(c *gin.Context) {
	var invoiceID = c.Param("id")

	_, err := invoiceCollection.DeleteOne(c.Request.Context(), bson.M{"invoice_id": invoiceID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Invoice deleted successfully", "status": http.StatusOK, "success": true})
}
