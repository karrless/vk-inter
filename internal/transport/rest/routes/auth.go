package routes

import (
	"context"

	"vk-inter/internal/transport/rest/controllers"
	"vk-inter/internal/transport/rest/interfaces"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(ctx *context.Context, r *gin.RouterGroup, authService interfaces.AuthService) {
	authController := controllers.NewAuthController(ctx, authService)
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authController.SignUp)
		authGroup.POST("/login", authController.LogIn)
	}
}
