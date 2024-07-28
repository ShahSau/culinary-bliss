package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
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

func GetOrderItems(c *gin.Context) {
	orders, err := orderItemCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer orders.Close(c.Request.Context())

	var allOrdersItems []bson.M

	for orders.Next(c.Request.Context()) {
		var orderItem bson.M
		if err = orders.Decode(&orderItem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		}
		allOrdersItems = append(allOrdersItems, orderItem)
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Order Items retrived successfully", "data": allOrdersItems, "status": http.StatusOK, "success": true})

}

func GetOrderItemsByOrder(c *gin.Context) {
	orderID := c.Param("id")
	allorders, err := ItemsByOrder(orderID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Order Items retrived successfully", "data": allorders, "status": http.StatusOK, "success": true})

}

func GetOrderItem(c *gin.Context) {
	var orderItemId = c.Param("id")
	var orderItem models.OrderItem

	defer c.Request.Body.Close()

	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errFetch := orderItemCollection.FindOne(c.Request.Context(), primitive.M{"order_item_id": orderItemId}).Decode(&orderItem)

	if errFetch != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errFetch.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Order Item created successfully", "data": orderItem, "status": http.StatusOK, "success": true})
}

func CreateOrderItem(c *gin.Context) {
	var orderItem models.OrderItem

	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var OrderItemPack OrderItemPack
	var order models.Order

	defer c.Request.Body.Close()

	order.Order_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Table_id = OrderItemPack.Table_id
	order_id := OrderItemOrderCreator(order)
	orderItems := []interface{}{}

	for _, item := range OrderItemPack.Order_items {
		item.Order_id = order_id
		item.ID = primitive.NewObjectID()
		item.Order_item_id = primitive.NewObjectID().Hex()
		item.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		item.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		item.Unit_price = toFixed(item.Unit_price, 2)
		orderItems = append(orderItems, item)
	}

	_, err := orderItemCollection.InsertMany(c.Request.Context(), orderItems)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": nil, "message": "Order Item created successfully", "data": orderItems, "status": http.StatusCreated, "success": true})

}

func UpdateOrderItem(c *gin.Context) {
	orderItemId := c.Param("id")
	var orderItem models.OrderItem

	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer c.Request.Body.Close()

	filter := bson.M{"order_item_id": orderItemId}

	var updateObj primitive.D

	if orderItem.Unit_price != 0 {
		updateObj = append(updateObj, bson.E{Key: "unit_price", Value: orderItem.Unit_price})
	}

	if quantity, err := strconv.Atoi(orderItem.Quantity); err == nil && quantity != 0 {
		updateObj = append(updateObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
	}

	if orderItem.Food_id != "" {
		updateObj = append(updateObj, bson.E{Key: "food_id", Value: orderItem.Food_id})
	}

	orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.UpdatedAt})

	err := orderItemCollection.FindOneAndUpdate(c.Request.Context(), filter, bson.D{{Key: "$set", Value: updateObj}}).Decode(&orderItem)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Order Item updated successfully", "data": orderItem, "status": http.StatusOK, "success": true})
}

func DeleteOrderItem(c *gin.Context) {
	fmt.Println("DeleteOrderItem")
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	return
}
