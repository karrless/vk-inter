package interfaces

import (
	"context"
)

type ListingService interface {
	CreateListing(ctx *context.Context) error
	GetListings(ctx *context.Context) error
}
