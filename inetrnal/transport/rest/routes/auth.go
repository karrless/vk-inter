package routes

import (
	"context"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(ctx *context.Context, r *gin.RouterGroup, authService *service.AuthService) {
	authController := controllers.NewAuthController(ctx, authService)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authController.SignUp)
		authGroup.POST("/login", authController.LogIn)
	}
}
