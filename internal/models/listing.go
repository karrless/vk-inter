package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Listing struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Title       string             `bson:"title" json:"title" validate:"required,min=3,max=100"`
	Description string             `bson:"description" json:"description" validate:"required,min=10,max=5000"`
	ImageURL    string             `bson:"image_url" json:"image_url" validate:"required,url,max=500"`
	Price       float64            `bson:"price" json:"price" validate:"required,gte=0,lte=1000000000"`
	OwnerID     primitive.ObjectID `bson:"owner_id" json:"owner_id"`
	OwnerLogin  string             `bson:"owner_login" json:"owner_login"`
	IsMyListing *bool              `bson:"is_my_listing,omitempty" json:"is_my_listing,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}
