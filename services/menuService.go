package services

import (
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

var menuCollection *mongo.Collection = database.GetCollection(database.DB, "menu")

func GetMenus(c *gin.Context) (models.ResponseMenu, error) {
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
				{Key: "name", Value: 1},
				{Key: "description", Value: 1},
				{Key: "start_date", Value: 1},
				{Key: "end_date", Value: 1},
				{Key: "menu_id", Value: 1},
			},
		},
	}

	record := bson.D{{Key: "$skip", Value: recordPerPage * (page - 1)}}
	limit := bson.D{{Key: "$limit", Value: recordPerPage}}

	result, errAgg := menuCollection.Aggregate(c.Request.Context(), mongo.Pipeline{matchStage, projectStage, record, limit})

	if errAgg != nil {
		return models.ResponseMenu{}, errAgg
	}

	var allMenus []bson.M
	if err = result.All(c.Request.Context(), &allMenus); err != nil {
		log.Fatal(err)
	}

	var menus []models.Menu
	for _, menu := range allMenus {
		var m models.Menu
		bsonBytes, _ := bson.Marshal(menu)
		bson.Unmarshal(bsonBytes, &m)
		menus = append(menus, m)
	}

	response := models.ResponseMenu{
		AllMenus:      menus,
		Page:          page,
		RecordPerPage: recordPerPage,
		StartIndex:    startIndex,
	}
	return response, nil
}

func GetMenuByID(id string, c *gin.Context) (models.Menu, error) {
	var menu models.Menu
	menuID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return menu, errors.New("invalid menu ID")
	}

	err = menuCollection.FindOne(c.Request.Context(), bson.M{"_id": menuID}).Decode(&menu)
	if err != nil {
		return menu, errors.New("menu not found")
	}

	return menu, nil
}

func CreateMenu(menu models.Menu, c *gin.Context) (models.Menu, error) {

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.Menu{}, errors.New("user is not admin")
	}

	var reqMenu models.Menu

	reqMenu.Name = menu.Name
	reqMenu.Description = menu.Description
	reqMenu.Start_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.End_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.ID = primitive.NewObjectID()
	reqMenu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.Menu_id = reqMenu.ID.Hex()

	_, err := menuCollection.InsertOne(c.Request.Context(), reqMenu)
	if err != nil {
		return models.Menu{}, err
	}

	return reqMenu, nil
}

func UpdateMenu(id string, menu models.Menu, c *gin.Context) (models.Menu, error) {

	var reqMenu models.Menu
	menuID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return reqMenu, errors.New("invalid menu ID")
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return reqMenu, errors.New("user is not admin")
	}

	reqMenu.Name = menu.Name
	reqMenu.Description = menu.Description
	reqMenu.Start_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.End_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	reqMenu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err = menuCollection.UpdateOne(c.Request.Context(), bson.M{"_id": menuID}, bson.M{"$set": reqMenu})
	if err != nil {
		return reqMenu, err
	}

	return reqMenu, nil
}

func DeleteMenu(id string, c *gin.Context) error {
	menuID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid menu ID")
	}

	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return errors.New("user is not admin")
	}

	_, err = menuCollection.DeleteOne(c.Request.Context(), bson.M{"_id": menuID})
	if err != nil {
		return err
	}

	return nil
}
