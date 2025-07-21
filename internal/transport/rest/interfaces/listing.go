package interfaces

import (
	"context"
	"vk-inter/internal/models"
)

type ListingService interface {
	CreateListing(ctx context.Context, title, description, imageURL string, price int, user *models.User) (*models.Listing, error)
	GetListings(ctx context.Context) error
}
