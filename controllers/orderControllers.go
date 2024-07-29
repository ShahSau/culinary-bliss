package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

func GetOrders(c *gin.Context) {
	orders, err := orderCollection.Find(c.Request.Context(), bson.M{}, nil)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer orders.Close(c.Request.Context())

	var results []models.Order

	for orders.Next(c.Request.Context()) {
		var order models.Order
		if err = orders.Decode(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": err.Error()})
		}
		results = append(results, order)
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Order retrived successfully", "data": results, "status": http.StatusOK, "success": true})
}

func GetOrder(c *gin.Context) {
	order_id := c.Param("id")

	if err := c.ShouldBindJSON(&order_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order

	defer c.Request.Body.Close()

	err := orderCollection.FindOne(c.Request.Context(), bson.M{"order_id": order_id}).Decode(&order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Order retrived successfully", "data": order, "status": http.StatusOK, "success": true})

}

func CreateOrder(c *gin.Context) {
	var orderReq models.Order
	var table models.Table

	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if orderReq.Table_id == "" {
		err := tableCollection.FindOne(c.Request.Context(), bson.M{"table_id": orderReq.Table_id}).Decode(&table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	orderReq.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderReq.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	orderReq.ID = primitive.NewObjectID()
	orderReq.Order_id = orderReq.ID.Hex()

	_, err := orderCollection.InsertOne(c.Request.Context(), orderReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": nil, "message": "Order created successfully", "data": orderReq, "status": http.StatusCreated, "success": true})
}

func UpdateOrder(c *gin.Context) {
	var reqOrder models.Order
	var Table models.Table

	var updateObj primitive.D

	orderId := c.Param("id")
	if err := c.ShouldBindJSON(&reqOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if reqOrder.Table_id != "" {
		err := menuCollection.FindOne(c.Request.Context(), bson.M{"table_id": reqOrder.Table_id}).Decode(&Table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		updateObj = append(updateObj, bson.E{Key: "table_id", Value: reqOrder.Table_id})
	}

	reqOrder.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err := orderCollection.UpdateOne(c.Request.Context(), bson.M{"order_id": orderId}, bson.D{{Key: "$set", Value: updateObj}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": nil, "message": "Order updated successfully", "status": http.StatusOK, "success": true})

}

func DeleteOrder(c *gin.Context) {
	fmt.Println("DeleteOrder")
}

func OrderItemOrderCreator(order models.Order) string {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	defer cancel()

	return order.Order_id
}
