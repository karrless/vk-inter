package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Login     string             `bson:"login"`
	Password  string             `bson:"hashed_password"`
	CreatedAt time.Time          `bson:"created_at"`
}
