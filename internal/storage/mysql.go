package storage

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"car-export-go/internal/config"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
