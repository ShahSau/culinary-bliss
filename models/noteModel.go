package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     *string            `json:"title" binding:"required" bson:"title"`
	Text      *string            `json:"text" binding:"required" bson:"text"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Note_id   string             `json:"note_id" binding:"required" bson:"note_id"`
}
