package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Restaurant struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Restaurant_id string             `json:"restaurant_id" binding:"required" bson:"restaurant_id"`
	Title         string             `json:"title" binding:"required" bson:"title"`
	Image         string             `json:"image" binding:"required" bson:"image"`
	Time          string             `json:"time" binding:"required" bson:"time"`
	Pickup        bool               `json:"pickup" binding:"required" bson:"pickup"`
	Delivery      bool               `json:"delivery" binding:"required" bson:"delivery"`
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt     time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Rating        float64            `json:"rating" binding:"required" bson:"rating"`
	RatingCount   int                `json:"ratingCount" binding:"required" bson:"ratingCount"`
	Menu          []Menu             `json:"menu" binding:"required" bson:"menu"`
}
type ResponseRestaurant struct {
	AllRestaurants []Restaurant `json:"all_restaurants"`
	Page           int          `json:"page"`
	RecordPerPage  int          `json:"record_per_page"`
	StartIndex     int          `json:"start_index"`
}
