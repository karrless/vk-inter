package rest

import (
	"context"
	"vk-inter/docs"
	"vk-inter/internal/transport/rest/interfaces"
	"vk-inter/internal/transport/rest/middlewares"
	"vk-inter/internal/transport/rest/routes"
	"vk-inter/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type RestConfig struct {
	Host string `env:"REST_HOST" env-default:"localhost"`
	Port string `env:"REST_PORT" env-default:"8080"`
}

type Server struct {
	ctx *context.Context
	cfg RestConfig
	r   *gin.Engine
}

func New(ctx *context.Context, cfg RestConfig, debug bool, secret string, authService interfaces.AuthService) *Server { //listingService interfaces.ListingService
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	r.Use(gin.Recovery())

	r.Use(middlewares.WithLogger(ctx))

	authMiddleware := middlewares.AuthMiddleware(secret)
	r.Use(authMiddleware)

	r.SetTrustedProxies([]string{"127.0.0.1", cfg.Host})
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = cfg.Host + ":" + cfg.Port
	docs.SwaggerInfo.Title = "VK-Inter API"
	docs.SwaggerInfo.Description = "API for auth and listings"
	docs.SwaggerInfo.Version = "0.1.0"

	routes.AuthRoutes(ctx, r.Group("/"), authService)
	// routes.ListingRoute(ctx, r.Group("/"), listingService)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{ctx: ctx, cfg: cfg, r: r}
}

func (s *Server) Run() error {
	logger.FromContext(*s.ctx).Info("Starting server", zap.String("port", s.cfg.Port))
	return s.r.Run(":" + s.cfg.Port)
}
