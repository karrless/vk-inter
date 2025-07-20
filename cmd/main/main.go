package main

import (
	"context"
	"log"
	"vk-inter/inetrnal/config"
	"vk-inter/pkg/logger"
)

func main() {
	ctx := context.Background()

	cfg := config.New()
	if cfg == nil {
		log.Fatal("Config is nil")
	}

	mainLogger := logger.New(cfg.Debug)
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)

	mainLogger.Info("Starting application")

}
