package repository

import (
	"context"
	"errors"
	"regexp"
	"time"
	"vk-inter/internal/models"
	"vk-inter/pkg/db/mongo"
	"vk-inter/pkg/errs"
	"vk-inter/pkg/logger"
	"vk-inter/pkg/utils"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuthRepo struct {
	*mongo.MongoDB
	collection mongoDriver.Collection
	validate   *validator.Validate
}

func NewAuthRepo(ctx context.Context, db *mongo.MongoDB) *AuthRepo {

	log := logger.FromContext(ctx)

	err := db.CreateUniqueIndex(ctx, "users", "login", true)
	if err != nil {
		log.Fatal("Failed to create index", zap.Error(err))
	}

	schema := bson.M{
		"bsonType": "object",
		"required": []string{"login", "hashed_password"},
		"properties": bson.M{
			"login": bson.M{
				"bsonType":    "string",
				"minLength":   3,
				"maxLength":   32,
				"pattern":     "^[\\p{L}\\p{N}_-]+$",
				"description": "must be 3-32  chars",
			},
			"hashed_password": bson.M{
				"bsonType": "string",
			},
		},
	}

	err = db.SetupValidation(ctx, schema, "users")
	if err != nil {
		log.Fatal("Failed to setup validation", zap.Error(err))
	}

	validate := validator.New()
	validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(fl.Param())
		return re.MatchString(fl.Field().String())
	})

	return &AuthRepo{
		MongoDB:    db,
		collection: *db.Collection("users"),
		validate:   validate,
	}
}

func (ar *AuthRepo) SignUp(ctx context.Context, user *models.User) (*models.User, error) {
	if err := ar.validate.Struct(user); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				switch fe.Field() {
				case "Login":
					return nil, errs.ErrInvalidLoginFormat
				case "Password":
					return nil, errs.ErrInvalidPasswordLength
				}
			}
		}
		return nil, err
	}

	var err error
	user.Password, err = utils.HashString(user.Password)
	if err != nil {
		return nil, err
	}
	user.CreatedAt = time.Now()

	res, err := ar.collection.InsertOne(ctx, user)
	if err != nil {
		if ar.MongoDB.IsDuplicateKeyError(err) {
			return nil, errs.ErrUserAlreadyExsist
		}
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	user.Password = ""
	return user, nil
}

func (ar *AuthRepo) CheckUser(ctx context.Context, user *models.User) (string, error) {
	if err := ar.validate.Struct(user); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				switch fe.Field() {
				case "Login":
					return "", errs.ErrUserNotFound
				case "Password":
					return "", errs.ErrWrongPassword
				}
			}
		}
		return "", err
	}
	userMongo, err := ar.GetByLogin(ctx, user.Login)
	if err != nil {
		return "", err
	}

	err = utils.CheckStirngHash(user.Password, userMongo.Password)
	if err != nil {
		return "", errs.ErrWrongPassword
	}

	return userMongo.ID.Hex(), nil
}

func (ar *AuthRepo) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	// Создаем фильтр для поиска по логину
	filter := bson.M{"login": login}

	var user models.User
	err := ar.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocuments) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (ar *AuthRepo) GetByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error) {
	filter := bson.M{"_id": userID}

	var user models.User
	err := ar.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongoDriver.ErrNoDocuments) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}
	user.Password = ""
	return &user, nil
}
