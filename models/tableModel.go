package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Number_of_guests *int               `json:"number_of_guests" binding:"required" bson:"number_of_guests"`
	Table_id         string             `json:"table_id" binding:"required" bson:"table_id"`
	Table_number     *int               `json:"table_number" binding:"required" bson:"table_number"`
	Table_status     *string            `json:"table_status" binding:"required" bson:"table_status"`
	CreatedAt        time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt        time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
