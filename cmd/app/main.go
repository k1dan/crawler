package main

import (
	"context"

	"github.com/k1dan/crawler/internal/app"
	"github.com/k1dan/crawler/internal/config"
	logger "github.com/k1dan/crawler/internal/logger"
)

const defaultLogLevel = "info"
const envFilePath = ".env"

func main() {
	ctx := context.Background()
	log := logger.New(defaultLogLevel)
	cfg := config.Load(log, envFilePath)

	application := app.New(log, cfg)
	application.Run(ctx, cfg.ParsingStartURL)
}
