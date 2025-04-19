package services

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/ShahSau/culinary-bliss/helpers"
	"github.com/ShahSau/culinary-bliss/models"
	"github.com/ShahSau/culinary-bliss/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers(c *gin.Context) (models.ResponseUser, error) {
	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.ResponseUser{}, errors.New("you are not authorized to view this resource")
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
				{Key: "first_name", Value: 1},
				{Key: "last_name", Value: 1},
				{Key: "email", Value: 1},
				{Key: "phone", Value: 1},
				{Key: "role", Value: 1},
			},
		},
	}

	record := bson.D{{Key: "$skip", Value: recordPerPage * (page - 1)}}
	limit := bson.D{{Key: "$limit", Value: recordPerPage}}

	result, errAgg := userCollection.Aggregate(c.Request.Context(), mongo.Pipeline{matchStage, projectStage, record, limit})

	if errAgg != nil {
		log.Fatal(errAgg)
		return models.ResponseUser{}, errAgg
	}

	var allUsers []bson.M
	if err = result.All(c.Request.Context(), &allUsers); err != nil {
		log.Fatal(err)
		return models.ResponseUser{}, err
	}

	var users []models.User
	for _, user := range allUsers {
		var userModel models.User
		bsonBytes, _ := bson.Marshal(user)
		bson.Unmarshal(bsonBytes, &userModel)
		users = append(users, userModel)
	}

	response := models.ResponseUser{
		AllUsers:      users,
		Page:          page,
		RecordPerPage: recordPerPage,
		StartIndex:    startIndex,
	}
	return response, nil
}

func GetUser(c *gin.Context, id string) (models.User, error) {
	userId, err := primitive.ObjectIDFromHex(id)

	var user models.User

	defer c.Request.Body.Close()

	err = userCollection.FindOne(c.Request.Context(), bson.D{{Key: "user_id", Value: userId}}).Decode(&user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func UpdateUser(c *gin.Context, id string, userReq models.User) (models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	defer c.Request.Body.Close()

	var updatedUser models.User

	err = userCollection.FindOne(c.Request.Context(), bson.M{"user_id": objectID}).Decode(&updatedUser)

	if err != nil {
		return models.User{}, err
	}

	updatedUser.Password = userReq.Password
	updatedUser.First_name = userReq.First_name
	updatedUser.Last_name = userReq.Last_name
	updatedUser.Email = userReq.Email
	updatedUser.ID = objectID
	updatedUser.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err = userCollection.UpdateOne(c.Request.Context(), bson.M{"user_id": objectID}, bson.D{{Key: "$set", Value: updatedUser}})
	if err != nil {
		return models.User{}, err
	}
	return updatedUser, nil
}

func DeleteUser(c *gin.Context, id string) (models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	defer c.Request.Body.Close()
	userEmail, _ := c.Get("first_name")
	var isAdmin = helpers.IsAdmin(userEmail.(string))

	if !isAdmin {
		return models.User{}, errors.New("unauthorized")
	}
	var deletedUser models.User

	_, err = userCollection.DeleteOne(c.Request.Context(), bson.M{"user_id": objectID})
	if err != nil {
		return models.User{}, err
	}

	return deletedUser, nil
}

func ResetPassword(c *gin.Context, userReq types.PasswordReset) (models.User, error) {
	var foundUser models.User
	err := userCollection.FindOne(c.Request.Context(), bson.D{{Key: "email", Value: userReq.Email}}).Decode(&foundUser)
	if err != nil {
		return models.User{}, err
	}
	passwordIsValid, msg := ComparePassword(foundUser.Password, userReq.OldPassword)

	if !passwordIsValid {
		return models.User{}, errors.New(msg)
	}

	new_password := HashPassword(userReq.NewPassword)
	foundUser.Password = new_password

	foundUser.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, err = userCollection.UpdateOne(c.Request.Context(), bson.D{{Key: "email", Value: userReq.Email}}, bson.D{{Key: "$set", Value: foundUser}})
	if err != nil {
		return models.User{}, err
	}
	return foundUser, nil
}
