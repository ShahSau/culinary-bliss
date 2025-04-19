package services

import (
	"context"
	"errors"
	"log"
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
var tableCollection *mongo.Collection = database.GetCollection(database.DB, "tables")

func GetOrders(c *gin.Context) (models.ResponseOrder, error) {
	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.ResponseOrder{}, errors.New("you are not authorized to view this resource")
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
			},
		},
	}
	record := bson.D{{Key: "$skip", Value: recordPerPage * (page - 1)}}
	limit := bson.D{{Key: "$limit", Value: recordPerPage}}

	result, errAgg := orderCollection.Aggregate(c.Request.Context(), mongo.Pipeline{matchStage, projectStage, record, limit})

	if errAgg != nil {
		return models.ResponseOrder{}, errAgg
	}

	var allOrders []bson.M
	if err = result.All(c.Request.Context(), &allOrders); err != nil {
		log.Fatal(err)
	}

	var orders []models.Order
	for _, order := range allOrders {
		var mappedOrder models.Order
		bsonBytes, _ := bson.Marshal(order)
		bson.Unmarshal(bsonBytes, &mappedOrder)
		orders = append(orders, mappedOrder)
	}

	response := models.ResponseOrder{
		AllOrders:     orders,
		Page:          page,
		RecordPerPage: recordPerPage,
		StartIndex:    startIndex,
	}
	return response, nil
}

func GetOrderById(c *gin.Context, orderId string) (models.Order, error) {
	var order models.Order

	err := orderCollection.FindOne(c.Request.Context(), bson.M{"order_id": orderId}).Decode(&order)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func CreateOrder(c *gin.Context, orderReq models.Order) (models.Order, error) {
	var table models.Table

	if err := c.ShouldBindJSON(&orderReq); err != nil {
		return models.Order{}, err
	}

	if orderReq.Table_id == "" {
		err := tableCollection.FindOne(c.Request.Context(), bson.M{"table_id": orderReq.Table_id}).Decode(&table)
		if err != nil {
			return models.Order{}, err
		}
	}
	orderReq.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderReq.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderReq.Order_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderReq.ID = primitive.NewObjectID()
	orderReq.Order_id = orderReq.ID.Hex()

	_, err := orderCollection.InsertOne(c.Request.Context(), orderReq)
	if err != nil {
		return models.Order{}, err
	}
	return orderReq, nil
}

func UpdateOrder(c *gin.Context, reqOrder models.Order, orderId string) (models.Order, error) {
	if err := c.ShouldBindJSON(&reqOrder); err != nil {
		return models.Order{}, err
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.Order{}, errors.New("user is not admin")
	}

	reqOrder.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	_, err := orderCollection.UpdateOne(c.Request.Context(), bson.M{"order_id": orderId}, bson.D{{Key: "$set", Value: reqOrder}})

	if err != nil {
		return models.Order{}, err
	}

	return reqOrder, nil
}

func DeleteOrder(c *gin.Context, orderId string) (models.Order, error) {
	var order models.Order

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.Order{}, errors.New("user is not admin")
	}

	err := orderCollection.FindOne(c.Request.Context(), bson.M{"order_id": orderId}).Decode(&order)
	if err != nil {
		return models.Order{}, err
	}

	_, err = orderCollection.DeleteOne(c.Request.Context(), bson.M{"order_id": orderId})
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
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
