package controllers

import (
	"context"
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

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order Items retrived successfully", "data": allOrdersItems, "status": http.StatusOK, "success": true})

}

func GetOrderItemsByOrder(c *gin.Context) {
	orderID := c.Param("id")
	allorders, err := ItemsByOrder(orderID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order Items retrived successfully", "data": allorders, "status": http.StatusOK, "success": true})

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

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order Item created successfully", "data": orderItem, "status": http.StatusOK, "success": true})
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

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Order Item created successfully", "data": orderItems, "status": http.StatusCreated, "success": true})

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

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order Item updated successfully", "data": orderItem, "status": http.StatusOK, "success": true})
}

func DeleteOrderItem(c *gin.Context) {
	orderItemId := c.Param("id")

	_, err := orderItemCollection.DeleteOne(c.Request.Context(), bson.M{"order_item_id": orderItemId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
