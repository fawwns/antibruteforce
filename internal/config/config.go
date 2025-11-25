package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	LimitLogin    int
	LimitPassword int
	LimitIP       int
}

func Load() *Config {
	// Загружаем .env, если есть (не обязательно).
	_ = godotenv.Load()

	cfg := &Config{}

	cfg.Port = getEnv("PORT", "8080")
	cfg.LimitLogin = getEnvAsInt("LIMIT_LOGIN", 10)
	cfg.LimitPassword = getEnvAsInt("LIMIT_PASSWORD", 100)
	cfg.LimitIP = getEnvAsInt("LIMIT_IP", 1000)

	log.Printf("Config loaded: %+v\n", cfg)
	return cfg
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
		log.Printf("warning: %s must be int, got '%s', using default %d", key, valueStr, defaultVal)
	}
	return defaultVal
}
