package controllers

import (
	"context"
	"vk-inter/internal/transport/rest/interfaces"

	"github.com/gin-gonic/gin"
)

type ListingController struct {
	ctx            *context.Context
	listingService interfaces.ListingService
}

func NewListingController(ctx *context.Context, listingService interfaces.ListingService) *ListingController {
	return &ListingController{
		ctx:            ctx,
		listingService: listingService,
	}
}

func (lc *ListingController) CreateListing(c *gin.Context) {

}

func (lc *ListingController) GetListings(c *gin.Context) {

}
