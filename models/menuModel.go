package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" binding:"required" bson:"name"`
	Description string             `json:"description" binding:"required" bson:"description"`
	Start_Date  *time.Time         `json:"start_date" binding:"required" bson:"start_date"`
	End_Date    *time.Time         `json:"end_date" binding:"required" bson:"end_date"`
	Menu_id     string             `json:"menu_id" binding:"required" bson:"menu_id"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
