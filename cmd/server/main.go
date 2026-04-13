// Package main is the entry point for the Wishlist API server.
//
// @title Wishlist API
// @version 1.0
// @description REST API сервис для создания вишлистов и бронирования подарков.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"log/slog"
	"os"
	"wishlist-api/internal/app"
	"wishlist-api/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	application, err := app.New(ctx, cfg)
	if err != nil {
		slog.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}

	application.Run()
}
