package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderItemPack struct {
	Table_id    string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.GetCollection(database.DB, "order_items")

func GetOrderItems(c *gin.Context) ([]models.OrderItem, error) {
	orders, err := orderItemCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		return nil, err
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return nil, errors.New("you are not authorized to view this resource")
	}

	defer orders.Close(c.Request.Context())

	var allOrdersItems []models.OrderItem

	for orders.Next(c.Request.Context()) {
		var orderItem models.OrderItem
		if err = orders.Decode(&orderItem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		}
		allOrdersItems = append(allOrdersItems, orderItem)
	}

	return allOrdersItems, nil
}

func GetOrderItemByID(id string, c *gin.Context) (models.OrderItem, error) {
	var orderItem models.OrderItem
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return orderItem, errors.New("invalid order item ID")
	}

	err = orderItemCollection.FindOne(c.Request.Context(), bson.M{"_id": objectID}).Decode(&orderItem)
	if err != nil {
		return orderItem, errors.New("order item not found")
	}

	return orderItem, nil
}

func CreateOrderItem(orderItem models.OrderItem, c *gin.Context) (models.OrderItem, error) {

	orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItem.ID = primitive.NewObjectID()
	orderItem.Order_item_id = orderItem.ID.Hex()

	_, err := orderItemCollection.InsertOne(c.Request.Context(), orderItem)

	if err != nil {
		return models.OrderItem{}, err
	}

	return orderItem, nil
}

func UpdateOrderItem(id string, updatedOrderItem models.OrderItem, c *gin.Context) (models.OrderItem, error) {
	orderItemID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return updatedOrderItem, errors.New("invalid order item ID")
	}

	err = orderItemCollection.FindOne(c.Request.Context(), bson.M{"order_item_id": updatedOrderItem}).Decode(&updatedOrderItem)

	if err != nil {
		return models.OrderItem{}, errors.New("order item not found")
	}
	var orderItem models.OrderItem
	orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItem.Quantity = updatedOrderItem.Quantity
	orderItem.Order_item_id = updatedOrderItem.Order_item_id
	orderItem.ID = updatedOrderItem.ID
	orderItem.CreatedAt = updatedOrderItem.CreatedAt
	orderItem.Order_id = updatedOrderItem.Order_id
	orderItem.Total_amount = updatedOrderItem.Total_amount
	_, err = orderItemCollection.UpdateOne(c.Request.Context(), bson.M{"_id": orderItemID}, bson.M{"$set": orderItem})
	if err != nil {
		return models.OrderItem{}, err
	}

	return orderItem, nil
}

func DeleteOrderItem(id string, c *gin.Context) (models.OrderItem, error) {
	orderItemID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.OrderItem{}, errors.New("invalid order item ID")
	}

	var orderItem models.OrderItem
	err = orderItemCollection.FindOne(c.Request.Context(), bson.M{"_id": orderItemID}).Decode(&orderItem)

	if err != nil {
		return models.OrderItem{}, errors.New("order item not found")
	}

	_, err = orderItemCollection.DeleteOne(c.Request.Context(), bson.M{"_id": orderItemID})
	if err != nil {
		return models.OrderItem{}, err
	}

	return orderItem, nil
}
