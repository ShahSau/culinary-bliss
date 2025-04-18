package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" binding:"required" bson:"name"`
	Description string             `json:"description" binding:"required" bson:"description"`
	Price       float64            `json:"price" binding:"required" bson:"price"`
	Image       string             `json:"image" binding:"required" bson:"image"`
	Food_id     string             `json:"food_id"  bson:"food_id"`
	Menu_id     string             `json:"menu_id" binding:"required" bson:"menu_id"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Response struct {
	AllFoods      []Food `json:"all_foods"`
	Page          int    `json:"page"`
	RecordPerPage int    `json:"record_per_page"`
	StartIndex    int    `json:"start_index"`
}
