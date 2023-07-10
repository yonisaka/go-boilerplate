package config

import (
	"os"
	"strconv"
)

type (
	Config struct {
		App      App
		MasterDB DB
		SlaveDB  DB
		Ws       Websocket
	}

	App struct {
		Name         string
		Env          string
		Port         int
		DefaultLang  string
		ReadTimeout  int
		WriteTimeout int
	}

	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		DB       string
	}

	Websocket struct {
		Host string
		Port int

		RoomIDLength      int
		MaxCachedMessages int
		MaxMessageLength  int
		Timeout           int
		RateLimitInterval int
		RateLimitMessages int
		MaxRooms          int
		MaxPeerPerRoom    int
		PeerHandleFormat  string
		RoomTimeout       int
		RoomAge           int
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
		MasterDB: DB{
			Host:     getEnv("POSTGRES_HOST_MASTER", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT_MASTER", 5432),
			User:     getEnv("POSTGRES_USER_MASTER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD_MASTER", "postgres"),
			DB:       getEnv("POSTGRES_DB_MASTER", "postgres"),
		},
		SlaveDB: DB{
			Host:     getEnv("POSTGRES_HOST_SLAVE", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT_SLAVE", 5432),
			User:     getEnv("POSTGRES_USER_SLAVE", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD_SLAVE", "postgres"),
			DB:       getEnv("POSTGRES_DB_SLAVE", "postgres"),
		},
		Ws: Websocket{
			Host:              getEnv("WS_HOST", "localhost"),
			Port:              getEnvAsInt("WS_PORT", 3000),
			RoomIDLength:      getEnvAsInt("WS_ROOM_ID_LENGTH", 8),
			MaxCachedMessages: getEnvAsInt("WS_MAX_CACHED_MESSAGES", 1000),
			MaxMessageLength:  getEnvAsInt("WS_MAX_MESSAGE_LENGTH", 3000),
			Timeout:           getEnvAsInt("WS_TIMEOUT", 3),
			RateLimitInterval: getEnvAsInt("WS_RATE_LIMIT_INTERVAL", 3),
			RateLimitMessages: getEnvAsInt("WS_RATE_LIMIT_MESSAGES", 25),
			MaxRooms:          getEnvAsInt("WS_MAX_ROOMS", 100),
			MaxPeerPerRoom:    getEnvAsInt("WS_MAX_PEER_PER_ROOM", 100),
			PeerHandleFormat:  getEnv("WS_PEER_HANDLE_FORMAT", "Guest-%d"),
			RoomTimeout:       getEnvAsInt("WS_ROOM_TIMEOUT", 10),
			RoomAge:           getEnvAsInt("WS_ROOM_AGE", 24),
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
