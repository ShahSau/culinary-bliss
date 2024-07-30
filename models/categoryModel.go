package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Category_id string             `json:"category_id,omitempty" bson:"category_id,omitempty"`
	Title       string             `json:"title,omitempty" binding:"required" bson:"title,omitempty"`
	Image       string             `json:"image,omitempty" bson:"image,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
