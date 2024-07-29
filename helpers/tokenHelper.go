package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ShahSau/culinary-bliss/database"
	"github.com/dgrijalva/jwt-go/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	User_id    string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.GetCollection(database.DB, "users")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, user_id string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_id:    user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Local().Add(time.Hour * time.Duration(24))), // 24 hours
		},
	}

	refreshClaims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		User_id:    user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Local().Add(time.Hour * time.Duration(24*7))), // 7 days
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, user_id string) {
	var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "$set", Value: bson.D{{Key: "token", Value: signedToken}, {Key: "refresh_token", Value: signedRefreshToken}}})
	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "$set", Value: bson.D{{Key: "updated_at", Value: updatedAt}}})

	_, err := userCollection.UpdateOne(ctx, bson.D{{Key: "user_id", Value: user_id}}, updateObj)

	if err != nil {
		log.Panic(err)
		return
	}

}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return nil, fmt.Errorf("expired token")
	}

	return claims, nil
}
