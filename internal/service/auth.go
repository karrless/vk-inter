package service

import (
	"context"
	"errors"
	"time"
	"vk-inter/internal/models"
	"vk-inter/pkg/errs"
	"vk-inter/pkg/jwt"
	"vk-inter/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthRepo interface {
	SignUp(ctx context.Context, user *models.User) (*models.User, error)
	CheckUser(ctx context.Context, user *models.User) (string, error)
	GetByID(ctx context.Context, userID primitive.ObjectID) (*models.User, error)
}

type AuthService struct {
	repo   AuthRepo
	secret string
}

func NewAuthService(repo AuthRepo, secret string) *AuthService {
	return &AuthService{
		repo:   repo,
		secret: secret,
	}
}

func (as *AuthService) SignUp(ctx context.Context, login, password string) (*models.User, error) {
	if err := utils.ValidatePassword(password, login); err != nil {
		return nil, err
	}
	return as.repo.SignUp(ctx, &models.User{
		Login:    login,
		Password: password,
	})
}

func (as *AuthService) LogIn(ctx context.Context, login, password string) (string, int, error) {
	user := &models.User{
		Login:    login,
		Password: password,
	}
	id, err := as.repo.CheckUser(ctx, user)
	if err != nil {
		if errors.Is(err, errs.ErrWrongPassword) || errors.Is(err, errs.ErrUserNotFound) {
			return "", 0, errs.ErrWrongPasswordOrLogin
		}
		return "", 0, err
	}
	duration := time.Hour * 48
	token, err := jwt.NewAccessToken(id, as.secret, duration)
	if err != nil {
		return "", 0, err
	}
	return token, int(duration.Seconds()), nil
}

func (as *AuthService) GetUserById(ctx context.Context, id string) (*models.User, error) {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}
	user, err := as.repo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
