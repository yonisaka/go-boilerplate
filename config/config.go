package config

import (
	"os"
	"strconv"
)

type (
	Config struct {
		App App
	}

	App struct {
		Name         string
		Env          string
		Port         int
		DefaultLang  string
		ReadTimeout  int
		WriteTimeout int
	}
)

func New() *Config {
	return &Config{
		App: App{
			Name:         getEnv("APP_NAME", "go-boilerplate"),
			Env:          getEnv("APP_ENV", "development"),
			Port:         getEnvAsInt("APP_PORT", 3000),
			DefaultLang:  getEnv("APP_DEFAULT_LANG", "en"),
			ReadTimeout:  getEnvAsInt("APP_READ_TIMEOUT", 10),
			WriteTimeout: getEnvAsInt("APP_WRITE_TIMEOUT", 10),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}

	if nextValue := os.Getenv(key); nextValue != "" {
		return nextValue
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return defaultVal
}
