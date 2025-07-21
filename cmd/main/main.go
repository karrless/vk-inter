// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 "Type 'Bearer {access_token}'"
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"vk-inter/internal/config"
	"vk-inter/internal/repository"
	"vk-inter/internal/service"
	rest "vk-inter/internal/transport"
	"vk-inter/pkg/db/mongo"
	"vk-inter/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg := config.New()
	if cfg == nil {
		log.Fatal("Config is nil")
	}
	ctx = context.WithValue(ctx, config.ConfigKey, cfg)

	mainLogger := logger.New(cfg.Debug)
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)

	if cfg.Secret == "test_key" {
		mainLogger.Warn("Secret key is default. Please, change it before run in production!")
	}

	db, err := mongo.New(ctx, cfg.MongoConfig)
	if err != nil {
		mainLogger.Fatal("Create MongoDB instanse error", zap.Error(err))
	}
	mainLogger.Debug("DB connected")

	authRepo := repository.NewAuthRepo(ctx, db)
	authService := service.NewAuthService(authRepo, cfg.Secret)

	restServer := rest.New(&ctx, cfg.RestConfig, cfg.Debug, cfg.Secret, authService)

	graceChannel := make(chan os.Signal, 1)
	signal.Notify(graceChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := restServer.Run(); err != nil {
			mainLogger.Fatal("failed to start server", zap.Error(err))
		}
	}()
	<-graceChannel
	db.Disconnect(ctx)

	mainLogger.Debug("MongoDB disconnected")
	mainLogger.Info("Graceful shutdown!")

}
