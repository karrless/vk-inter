package controllers

import (
	"context"
	"net/http"
	"vk-inter/internal/transport/rest/interfaces"
	"vk-inter/pkg/errs"

	"github.com/gin-gonic/gin"
)

type ListingController struct {
	ctx            *context.Context
	listingService interfaces.ListingService
	authService    interfaces.AuthService
}

func NewListingController(ctx *context.Context, listingService interfaces.ListingService, authService interfaces.AuthService) *ListingController {
	return &ListingController{
		ctx:            ctx,
		listingService: listingService,
		authService:    authService,
	}
}

type CreateListingRequest struct {
	Title       string
	Description string
	ImageURL    string
	Price       int
}

// @Summary	Create listing endpoint
// @Tags		listing
// @Security	BearerAuth
// @Accept		json
// @Produce	json
// @Param		request	body	CreateListingRequest	true	"Listing info"
// @Success	201		{object}	CreateListingResponse
// @Failure	400		{object}	ErrorResponse
// @Failure	403		{object}	ErrorResponse
// @Failure	409		{object}	ErrorResponse
// @Router		/listings [post]
func (lc *ListingController) CreateListing(c *gin.Context) {
	isAuth, exists := c.Get("isAuthenticated")
	tokenSubject, _ := c.Get("id")
	if !exists || !isAuth.(bool) || tokenSubject == nil {
		resp := ErrorResponse{
			Error:   errs.ErrUnauthorized.Error(),
			Message: "Please login before create listing",
		}
		c.JSON(http.StatusUnauthorized, resp)
		return
	}
	user, err := lc.authService.GetUserById(*lc.ctx, tokenSubject.(string))
	if err != nil {
		resp := ErrorResponse{
			Error:   errs.ErrUnauthorized.Error(),
			Message: "Please login before create listing",
		}
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	var req CreateListingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = lc.listingService.CreateListing(*lc.ctx, req.Title, req.Description, req.ImageURL, req.Price, user)
	if err != nil {
		status := http.StatusInternalServerError
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
}

func (lc *ListingController) GetListings(c *gin.Context) {

}
