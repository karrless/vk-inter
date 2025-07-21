package repository

import (
	"context"
	"errors"
	"math"
	"regexp"
	"time"
	"vk-inter/internal/models"
	"vk-inter/pkg/db/mongo"
	"vk-inter/pkg/errs"
	"vk-inter/pkg/logger"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type ListingRepo struct {
	*mongo.MongoDB
	collection mongoDriver.Collection
	validate   *validator.Validate
}

func NewListingRepo(ctx context.Context, db *mongo.MongoDB) *ListingRepo {
	log := logger.FromContext(ctx)

	err := db.CreateUniqueIndex(ctx, "listings", "title", false) // optional: title uniqueness
	if err != nil {
		log.Fatal("Failed to create index for listings", zap.Error(err))
	}

	schema := bson.M{
		"bsonType": "object",
		"required": []string{"title", "description", "image_url", "price", "owner_id", "owner_login"},
		"properties": bson.M{
			"title": bson.M{
				"bsonType":    "string",
				"minLength":   3,
				"maxLength":   100,
				"description": "must be 3-100 chars",
			},
			"description": bson.M{
				"bsonType":    "string",
				"minLength":   10,
				"maxLength":   5000,
				"description": "must be 10-5000 chars",
			},
			"image_url": bson.M{
				"bsonType":    "string",
				"maxLength":   500,
				"description": "must be valid URL",
			},
			"price": bson.M{
				"bsonType":    "double",
				"minimum":     0,
				"maximum":     1000000000,
				"description": "must be a decimal with up to 2 decimal places",
			},
			"owner_id": bson.M{
				"bsonType": "objectId",
			},
			"owner_login": bson.M{
				"bsonType": "string",
			},
		},
	}

	err = db.SetupValidation(ctx, schema, "listings")
	if err != nil {
		log.Fatal("Failed to setup validation for listings", zap.Error(err))
	}

	validate := validator.New()
	validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(fl.Param())
		return re.MatchString(fl.Field().String())
	})

	return &ListingRepo{
		MongoDB:    db,
		collection: *db.Collection("listings"),
		validate:   validate,
	}
}

func (lr *ListingRepo) CreateListing(ctx context.Context, listing *models.Listing) (*models.Listing, error) {
	if err := lr.validate.Struct(listing); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				switch fe.Field() {
				case "Title":
					return nil, errs.ErrListingInvalidTitle
				case "Description":
					return nil, errs.ErrListingInvalidDescription
				case "ImageURL":
					return nil, errs.ErrListingInvalidImageURL
				case "Price":
					return nil, errs.ErrListingInvalidPrice
				}
			}
		}
		return nil, err
	}

	listing.CreatedAt = time.Now()
	listing.Price = math.Round(listing.Price*100) / 100

	doc := bson.M{
		"title":       listing.Title,
		"description": listing.Description,
		"image_url":   listing.ImageURL,
		"price":       listing.Price,
		"owner_id":    listing.OwnerID,
		"owner_login": listing.OwnerLogin,
		"created_at":  listing.CreatedAt,
	}

	res, err := lr.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	listing.ID = res.InsertedID.(primitive.ObjectID)

	return listing, nil
}

func (lr *ListingRepo) GetListings(
	ctx context.Context,
	page, limit int,
	sortBy, order string,
	minPrice, maxPrice float64,
	currentUserID primitive.ObjectID,
) ([]*models.Listing, error) {
	skip := (page - 1) * limit

	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}

	filter := bson.M{
		"price": bson.M{
			"$gte": minPrice,
			"$lte": maxPrice,
		},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: sortOrder}}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := lr.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var listings []*models.Listing
	for cursor.Next(ctx) {
		var l models.Listing
		if err := cursor.Decode(&l); err != nil {
			continue
		}

		if currentUserID != primitive.NilObjectID {
			val := l.OwnerID == currentUserID
			l.IsMyListing = &val
		}

		listings = append(listings, &l)
	}

	return listings, nil
}
