package services

import (
	"errors"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var invoiceCollection *mongo.Collection = database.GetCollection(database.DB, "invoice")

func GetInvoices(c *gin.Context) ([]models.InvoiceViewFormat, error) {
	invoices, err := invoiceCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		return nil, err
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return nil, errors.New("you are not authorized to view this resource")
	}

	defer invoices.Close(c.Request.Context())

	var results []models.InvoiceViewFormat

	for invoices.Next(c.Request.Context()) {
		var invoice models.InvoiceViewFormat
		err := invoices.Decode(&invoice)
		if err != nil {
			return nil, err
		}
		results = append(results, invoice)
	}

	return results, nil
}

func GetInvoiceByID(c *gin.Context, invoiceID string) (models.InvoiceViewFormat, error) {
	var invoice models.InvoiceViewFormat

	id, err := primitive.ObjectIDFromHex(invoiceID)
	if err != nil {
		return models.InvoiceViewFormat{}, err
	}

	err = invoiceCollection.FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&invoice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.InvoiceViewFormat{}, errors.New("invoice not found")
		}
		return invoice, err
	}

	return invoice, nil
}

func CreateInvoice(c *gin.Context, reqInvoice models.Invoice) (models.Invoice, error) {
	var orderCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

	var order models.Order

	err := orderCollection.FindOne(c.Request.Context(), bson.M{"order_id": reqInvoice.Order_id}).Decode(&order)

	if err != nil {
		return models.Invoice{}, err
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
		return models.Invoice{}, err
	}
	return invoice, nil
}

func UpdateInvoice(c *gin.Context, invoiceID string, reqInvoice models.Invoice) (models.Invoice, error) {
	var invoice models.Invoice

	id, err := primitive.ObjectIDFromHex(invoiceID)
	if err != nil {
		return models.Invoice{}, err
	}

	err = invoiceCollection.FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&invoice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Invoice{}, errors.New("invoice not found")
		}
		return invoice, err
	}

	var orderCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

	var order models.Order

	err = orderCollection.FindOne(c.Request.Context(), bson.M{"order_id": invoice.Order_id}).Decode(&order)

	if err != nil {
		return models.Invoice{}, errors.New(err.Error() + " Order not found")
	}

	var updateObj models.Invoice

	updateObj.Order_id = reqInvoice.Order_id
	updateObj.Payment_method = reqInvoice.Payment_method
	updateObj.Payment_status = "PENDING"
	updateObj.Payment_due_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj.ID = invoice.ID
	updateObj.Invoice_id = invoice.Invoice_id
	updateObj.Total_amount = order.Total_amount

	_, err = invoiceCollection.UpdateOne(c.Request.Context(), bson.M{"invoice_id": invoiceID}, bson.D{{Key: "$set", Value: updateObj}})

	if err != nil {
		return models.Invoice{}, err
	}
	return updateObj, nil
}

func DeleteInvoice(c *gin.Context, invoiceID string) error {
	id, err := primitive.ObjectIDFromHex(invoiceID)
	if err != nil {
		return err
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return errors.New("you are not authorized to delete this resource")
	}
	_, err = invoiceCollection.DeleteOne(c.Request.Context(), bson.M{"_id": id})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("invoice not found")
		}
		return err
	}

	return nil
}
