package di

import (
	"github.com/joho/godotenv"
	"github.com/yonisaka/go-boilerplate/config"
	"github.com/yonisaka/go-boilerplate/pkg/logger"
)

func GetConfig() *config.Config {
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file")
	}

	return config.New()
}
