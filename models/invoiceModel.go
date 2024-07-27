package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Invoice_id       string             `json:"invoice_id" binding:"required" bson:"invoice_id"`
	Order_id         string             `json:"order_id" binding:"required" bson:"order_id"`
	Payment_method   string             `json:"payment_method" binding:"required" validate:"eq=CARD|eq=CASH|eq=" bson:"payment_method"`
	Payment_status   string             `json:"payment_status" binding:"required" bson:"payment_status"`
	Payment_due_date time.Time          `json:"payment_due_date" binding:"required" bson:"payment_due_date"`
	Total_amount     float64            `json:"total_amount" binding:"required" bson:"total_amount"`
	CreatedAt        time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt        time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
