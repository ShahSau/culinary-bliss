package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection = database.GetCollection(database.DB, "orders")

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
	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}
	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	page, errp := strconv.Atoi(c.Query("page"))
	if errp != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(c.Query("startIndex"))

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	projectStage := bson.D{
		{
			Key: "$project", Value: bson.D{
				{Key: "id", Value: 1},
				{Key: "table_id", Value: 1},
				{Key: "order_status", Value: 1},
				{Key: "order_date", Value: 1},
				{Key: "total_amount", Value: 1},
				{Key: "order_id", Value: 1},
			},
		},
	}

	record := bson.D{{Key: "$skip", Value: recordPerPage * (page - 1)}}
	limit := bson.D{{Key: "$limit", Value: recordPerPage}}

	result, errAgg := orderCollection.Aggregate(c.Request.Context(), mongo.Pipeline{matchStage, projectStage, record, limit})

	if errAgg != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errAgg.Error()})
		return
	}

	var allOrders []bson.M
	if err = result.All(c.Request.Context(), &allOrders); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": allOrders, "page": page, "recordPerPage": recordPerPage, "startIndex": startIndex})
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

	var order models.Order

	defer c.Request.Body.Close()

	err := orderCollection.FindOne(c.Request.Context(), bson.M{"order_id": order_id}).Decode(&order)

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
	orderReq.Order_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderReq.ID = primitive.NewObjectID()
	orderReq.Order_id = orderReq.ID.Hex()

	_, err := orderCollection.InsertOne(c.Request.Context(), orderReq)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Order created successfully", "data": orderReq, "status": http.StatusCreated, "success": true})
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
	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to view this resource"})
		return
	}

	reqOrder.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err := orderCollection.UpdateOne(c.Request.Context(), bson.M{"order_id": orderId}, bson.D{{Key: "$set", Value: reqOrder}})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order updated successfully", "status": http.StatusOK, "success": true, "data": reqOrder})

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

	_, err := orderCollection.DeleteOne(c.Request.Context(), bson.M{"order_id": orderId})

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

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Order deleted successfully", "status": http.StatusOK, "success": true, "data": nil})
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
