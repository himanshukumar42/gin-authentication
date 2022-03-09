package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterInput struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Username string             `json:"username,omitempty" validate:"required" binding:"required"`
	Password string             `json:"password,omitempty" validate:"required" binding:"required"`
}
