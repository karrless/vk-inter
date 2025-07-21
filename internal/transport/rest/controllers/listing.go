package controllers

import (
	"context"
	"errors"
	"math"
	"net/http"
	"time"
	"vk-inter/internal/transport/rest/interfaces"
	"vk-inter/pkg/errs"
	"vk-inter/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
}
type CreateListingResponse struct {
	ID          primitive.ObjectID `json:"_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	ImageURL    string             `json:"image_url"`
	Price       float64            `json:"price"`
	OwnerID     primitive.ObjectID `json:"owner_id"`
	OwnerLogin  string             `json:"owner_login"`
	CreatedAt   time.Time          `json:"created_at"`
}

// @Summary	Create listing endpoint
// @Tags		listing
// @Security	BearerAuth
// @Accept		json
// @Produce	json
// @Param		request	body	CreateListingRequest	true	"Listing info"
// @Success	201		{object}	CreateListingResponse
// @Failure	400		{object}	ErrorResponse
// @Failure	401		{object}	ErrorResponse
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

	listing, err := lc.listingService.CreateListing(*lc.ctx, req.Title, req.Description, req.ImageURL, req.Price, user)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, errs.ErrListingInvalidTitle) ||
			errors.Is(err, errs.ErrListingInvalidDescription) ||
			errors.Is(err, errs.ErrListingInvalidImageURL) ||
			errors.Is(err, errs.ErrListingInvalidPrice) {
			status = http.StatusBadRequest
		}
		if _, ok := err.(utils.ImageError); ok {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	resp := CreateListingResponse{
		ID:          listing.ID,
		Title:       listing.Title,
		Description: listing.Description,
		ImageURL:    listing.ImageURL,
		Price:       listing.Price,
		OwnerID:     listing.OwnerID,
		OwnerLogin:  listing.OwnerLogin,
		CreatedAt:   listing.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

// @Summary      Get listings
// @Tags         listing
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "Page number (default: 1)"
// @Param        limit       query     int     false  "Items per page (default: 10)"
// @Param        sort_by     query     string  false  "Sort field: created_at or price"
// @Param        order       query     string  false  "Sort order: asc or desc"
// @Param        min_price   query     number  false  "Minimum price filter"
// @Param        max_price   query     number  false  "Maximum price filter"
// @Success      200         {array}   models.Listing
// @Failure      401         {object}  ErrorResponse
// @Failure      500         {object}  ErrorResponse
// @Router       /listings [get]
func (lc *ListingController) GetListings(c *gin.Context) {
	var currentUserID primitive.ObjectID = primitive.NilObjectID

	if isAuth, exists := c.Get("isAuthenticated"); exists && isAuth.(bool) {
		if id, ok := c.Get("id"); ok {
			user, err := lc.authService.GetUserById(*lc.ctx, id.(string))
			if err == nil {
				currentUserID = user.ID
			}
		}
	}

	page := utils.ParseQueryInt(c, "page", 1)
	limit := utils.ParseQueryInt(c, "limit", 10)
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := c.DefaultQuery("order", "desc")
	minPrice := utils.ParseQueryFloat(c, "min_price", 0)
	maxPrice := utils.ParseQueryFloat(c, "max_price", math.MaxFloat64)

	listings, err := lc.listingService.GetListings(*lc.ctx, page, limit, sortBy, order, minPrice, maxPrice, currentUserID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, errs.ErrPriceSorting) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, listings)
}
