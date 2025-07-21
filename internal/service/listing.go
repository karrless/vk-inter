package service

import (
	"context"
	"math"

	"vk-inter/internal/models"
	"vk-inter/pkg/errs"
	"vk-inter/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListingRepo interface {
	CreateListing(ctx context.Context, listing *models.Listing) (*models.Listing, error)
	GetListings(ctx context.Context, page, limit int, sortBy, order string, minPrice, maxPrice float64, currentUserID primitive.ObjectID) ([]*models.Listing, error)
}

type ListingService struct {
	repo ListingRepo
}

func NewListingService(repo ListingRepo) *ListingService {
	return &ListingService{
		repo: repo,
	}
}

func (ls *ListingService) CreateListing(ctx context.Context, title, description, imageURL string, price float64, user *models.User) (*models.Listing, error) {
	err := utils.ValidateImageURL(imageURL)
	if err != nil {
		return nil, err
	}
	listing := models.Listing{
		Title:       title,
		Description: description,
		ImageURL:    imageURL,
		Price:       price,
		OwnerID:     user.ID,
		OwnerLogin:  user.Login,
	}
	return ls.repo.CreateListing(ctx, &listing)
}

func (ls *ListingService) GetListings(
	ctx context.Context,
	page, limit int,
	sortBy, order string,
	minPrice, maxPrice float64,
	currentUserID primitive.ObjectID,
) ([]*models.Listing, error) {
	// Ограничения по лимиту и странице
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Проверка допустимых значений сортировки
	if sortBy != "created_at" && sortBy != "price" {
		sortBy = "created_at"
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Безопасная фильтрация цены
	if minPrice < 0 {
		minPrice = 0
	}
	if maxPrice <= 0 || maxPrice > 1_000_000_000 {
		maxPrice = math.MaxFloat64
	}
	if maxPrice < minPrice {
		return []*models.Listing{}, errs.ErrPriceSorting
	}

	return ls.repo.GetListings(ctx, page, limit, sortBy, order, minPrice, maxPrice, currentUserID)
}
