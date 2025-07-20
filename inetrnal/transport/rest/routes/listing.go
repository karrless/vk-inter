package routes

import (
	"context"

	"github.com/gin-gonic/gin"
)

func ListingRoute(ctx *context.Context, r *gin.RouterGroup, listingService *service.ListingService) {
	listingController := controllers.NewListingController(ctx, listingService)
	authGroup := r.Group("/listings")
	{
		authGroup.GET("/", listingController.CreateListing)
		authGroup.POST("/", listingController.GetListings)
	}
}
