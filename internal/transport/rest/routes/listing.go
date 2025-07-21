package routes

import (
	"context"
	"vk-inter/internal/transport/rest/controllers"
	"vk-inter/internal/transport/rest/interfaces"

	"github.com/gin-gonic/gin"
)

func ListingRoute(ctx *context.Context, r *gin.RouterGroup, listingService interfaces.ListingService, authService interfaces.AuthService) {
	listingController := controllers.NewListingController(ctx, listingService, authService)
	authGroup := r.Group("/listings")
	{
		authGroup.GET("/", listingController.CreateListing)
		authGroup.POST("/", listingController.GetListings)
	}
}
