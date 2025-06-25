// Package config содержит работу с конфигурацией и .env переменными
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит настройки подключения к БД и экспорта
type Config struct {
	DBHost    string
	DBPort    string
	DBName    string
	DBUser    string
	DBPass    string
	ExportDir string
}

// LoadConfig загружает конфигурацию из .env и переменных окружения
func LoadConfig() *Config {
	_ = godotenv.Load(".env")
	cfg := &Config{
		DBHost:    getEnv("DB_HOST", "127.0.0.1"),
		DBPort:    getEnv("DB_PORT", "3306"),
		DBName:    getEnv("DB_DATABASE", ""),
		DBUser:    getEnv("DB_USERNAME", ""),
		DBPass:    getEnv("DB_PASSWORD", ""),
		ExportDir: getEnv("EXPORT_DIR", "exports"),
	}
	if cfg.DBName == "" || cfg.DBUser == "" {
		log.Fatal("DB_DATABASE and DB_USERNAME are required in .env")
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
