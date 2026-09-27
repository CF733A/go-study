package config

import (
	"log"
	"os"
)

type Config struct {
	Port         string
	DadataApiKey string
	DadataSecret string
	JWTsecret    string
}

func LoadConfig() *Config {
	config := &Config{
		Port:         getEnv("PORT", "8080"),
		DadataApiKey: getEnv("DADATA_API_KEY", ""),
		DadataSecret: getEnv("DADATA_SECRET_KEY", ""),
		JWTsecret:    getEnv("JWT_SECRET", ""),
	}

	if config.JWTsecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	log.Println("config loaded")
	return config
}

func getEnv(key, defValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defValue
}
