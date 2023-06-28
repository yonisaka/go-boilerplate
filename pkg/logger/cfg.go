package logger

type Config struct {
	Debug       bool   `json:"debug"`
	Environment string `json:"environment"`
	Level       string `json:"level"`
	ServiceName string `json:"service_name"`
}
