package interfaces

import (
	"context"
	"vk-inter/internal/models"
)

type AuthService interface {
	LogIn(ctx context.Context, login, password string) (string, int, error)
	SignUp(ctx context.Context, login, password string) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
}
