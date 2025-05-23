package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Order_id     string             `json:"order_id"  bson:"order_id"`
	Table_id     string             `json:"table_id" binding:"required" bson:"table_id"`
	Order_status string             `json:"order_status" binding:"required" bson:"order_status"`
	Order_date   time.Time          `json:"order_date" bson:"order_date"`
	Total_amount float64            `json:"total_amount" binding:"required" bson:"total_amount"`
	CreatedAt    time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type ResponseOrder struct {
	AllOrders     []Order `json:"all_orders"`
	Page          int     `json:"page"`
	RecordPerPage int     `json:"record_per_page"`
	StartIndex    int     `json:"start_index"`
}
