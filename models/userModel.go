package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	First_name string             `json:"first_name" binding:"required" bson:"first_name"`
	Last_name  string             `json:"last_name" binding:"required" bson:"last_name"`
	Email      string             `json:"email" binding:"required" bson:"email"`
	Password   string             `json:"password" binding:"required" bson:"password"`
	Avatar     string             `json:"avatar" bson:"avatar"`
	Phone      string             `json:"phone" bson:"phone"`
	Role       string             `json:"role" binding:"required" validate:"eq=Admin|eq=User" bson:"role"`
	CreatedAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	User_id    string             `json:"user_id" binding:"required" bson:"user_id"`
}
