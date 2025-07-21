package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Listing struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title" validate:"required,min=3,max=100"`
	Description string             `bson:"description" validate:"required,min=10,max=5000"`
	ImageURL    string             `bson:"image_url" validate:"required,url,max=500"`
	Price       int                `bson:"price" validate:"required,min=0,max=1000000000"`
	OwnerID     primitive.ObjectID `bson:"owner_id"`
	OwnerLogin  string             `bson:"owner_login"`
	IsImOwner   bool               `bson:"is_im_owner"`
	CreatedAt   time.Time          `bson:"created_at"`
}
