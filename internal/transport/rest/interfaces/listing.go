package interfaces

import (
	"context"
	"vk-inter/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListingService interface {
	CreateListing(ctx context.Context, title, description, imageURL string, price float64, user *models.User) (*models.Listing, error)
	GetListings(ctx context.Context, page, limit int, sortBy, order string, minPrice, maxPrice float64, currentUserID primitive.ObjectID) ([]*models.Listing, error)
}
