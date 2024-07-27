package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Food_id       string             `json:"food_id" binding:"required" bson:"food_id"`
	Order_id      string             `json:"order_id" binding:"required" bson:"order_id"`
	Order_item_id string             `json:"order_item_id" binding:"required" bson:"order_item_id"`
	Quantity      string             `json:"quantity" binding:"required" validate:"eq=S|eq=M|eq=L" bson:"quantity"`
	Unit_price    float64            `json:"price" binding:"required" bson:"unit_price"`
	Total_amount  float64            `json:"total_amount" binding:"required" bson:"total_amount"`
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt     time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
