package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Login     string             `bson:"login" validate:"required,min=3,max=32,regexp=^[\\p{L}\\p{N}_-]+$"`
	Password  string             `bson:"hashed_password" validate:"required,min=8"`
	CreatedAt time.Time          `bson:"created_at"`
}
