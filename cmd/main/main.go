package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"vk-inter/inetrnal/config"
	"vk-inter/pkg/db/mongo"
	"vk-inter/pkg/logger"
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
	mainLogger.Debug("Logger init")

	db, err := mongo.New(ctx, cfg.MongoConfig)
	if err != nil {
		mainLogger.Fatal(err.Error())
	}
	mainLogger.Debug("DB connected")

	graceChannel := make(chan os.Signal, 1)
	signal.Notify(graceChannel, syscall.SIGINT, syscall.SIGTERM)

	<-graceChannel
	db.Disconnect(ctx)

	mainLogger.Debug("MongoDB disconnected")
	mainLogger.Info("Graceful shutdown!")

}
