package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/services"
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

// @Summary Get Order Items
// @Description Get Order Items
// @Tags Admin
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /orderItems [get]
func GetOrderItems(c *gin.Context) {
	allOrdersItems, err := services.GetOrderItems(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order Items retrived successfully", "data": allOrdersItems, "status": http.StatusOK, "success": true})
}

// @Summary Get Order Item
// @Description Get Order Item
// @Tags User
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Order Item ID"
// @Success 200 {object} models.OrderItem
// @Failure 400 {object} string
// @Router /orderItem/{id} [get]
func GetOrderItem(c *gin.Context) {
	var orderItemId = c.Param("id")

	orderItem, err := services.GetOrderItemByID(orderItemId, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order Item created successfully", "data": orderItem, "status": http.StatusOK, "success": true})
}

// @Summary Create Order Item
// @Description Create Order Item
// @Tags User
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param orderItem body models.OrderItem true "Order Item Object"
// @Success 201 {object} models.OrderItem
// @Failure 400 {object} string
// @Router /orderItem [post]
func CreateOrderItem(c *gin.Context) {
	var orderItem models.OrderItem
	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderItem, err := services.CreateOrderItem(orderItem, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Order Item created successfully", "data": orderItem, "status": http.StatusCreated, "success": true})
}

// @Summary Update Order Item
// @Description Update Order Item
// @Tags User
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Order Item ID"
// @Param orderItem body models.OrderItem true "Order Item Object"
// @Success 200 {object} models.OrderItem
// @Failure 400 {object} string
// @Router /orderItem/{id} [put]
func UpdateOrderItem(c *gin.Context) {
	orderItemId := c.Param("id")
	var reqorderItem models.OrderItem

	if err := c.ShouldBindJSON(&reqorderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderItem, err := services.UpdateOrderItem(orderItemId, reqorderItem, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order Item updated successfully", "data": orderItem, "status": http.StatusOK, "success": true})
}

// @Summary Delete Order Item
// @Description Delete Order Item
// @Tags User
// @Accept json
// @Produce json
// @Security		BearerAuth
// @param Authorization header string true "Token"
// @Param id path string true "Order Item ID"
// @Success 200 {string} string	"Order Item deleted successfully"
// @Failure 400 {object} string
// @Router /orderItem/{id} [delete]
func DeleteOrderItem(c *gin.Context) {
	orderItemId := c.Param("id")

	_, err := services.GetOrderItemByID(orderItemId, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": fmt.Sprintf("Order Item with ID %s deleted successfully", orderItemId), "status": http.StatusOK, "success": true})
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {

	var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "food"}, {Key: "localField", Value: "food_id"}, {Key: "foreignField", Value: "food_id"}, {Key: "as", Value: "food"}}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$food"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	lookupOrderStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "order"}, {Key: "localField", Value: "order_id"}, {Key: "foreignField", Value: "order_id"}, {Key: "as", Value: "order"}}}}
	unwindOrderStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$order"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	lookupTableStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "table"}, {Key: "localField", Value: "order.table_id"}, {Key: "foreignField", Value: "table_id"}, {Key: "as", Value: "table"}}}}
	unwindTableStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$table"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "amount", Value: "$food.price"},
			{Key: "total_count", Value: 1},
			{Key: "food_name", Value: "$food.food_name"},
			{Key: "food_image", Value: "$food.food_image"},
			{Key: "table_number", Value: "$table.table_number"},
			{Key: "table_id", Value: "$table.table_id"},
			{Key: "order_id", Value: "$order.order_id"},
			{Key: "price", Value: "$food.price"},
			{Key: "quantity", Value: 1},
		}},
	}

	groupStage := bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"food_name", "$food_name"}, {"food_image", "$food_image"}, {"table_number", "$table_number"}, {"price", "$price"}, {"quantity", "$quantity"}, {"total_count", bson.D{{"$sum", 1}}}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"order_items", bson.D{{"$push", "$$ROOT"}}}}

	projectStage2 := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "id", Value: 0},
			{Key: "payment_due", Value: 1},
			{Key: "total_count", Value: 1},
			{Key: "table_number", Value: "$_id.table_number"},
			{Key: "order_items", Value: 1},
		}},
	}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, unwindStage, lookupOrderStage, unwindOrderStage, lookupTableStage, unwindTableStage, projectStage, groupStage, projectStage2})

	if err != nil {
		return nil, err
	}

	if err = result.All(ctx, &OrderItems); err != nil {
		return nil, err
	}

	return OrderItems, nil
}
