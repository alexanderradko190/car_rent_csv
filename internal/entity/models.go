// Package entity содержит бизнес-сущности (модели) проекта
package entity

import (
	"database/sql"
	"time"
)

// Car описывает автомобиль
type Car struct {
	ID              uint `gorm:"primaryKey"`
	Make            string
	Model           string
	Year            int
	Vin             string
	LicensePlate    string
	CarClass        string
	Power           int
	HourlyRate      float64
	Status          string
	CurrentRenterID sql.NullInt64
}

// Client описывает клиента
type Client struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string
}

// RentHistory описывает историю аренды автомобиля
type RentHistory struct {
	ID        uint `gorm:"primaryKey"`
	ClientID  uint
	CarID     uint
	StartTime time.Time
	EndTime   time.Time
	TotalCost float64
}

// RentalRequest описывает заявку на аренду автомобиля
type RentalRequest struct {
	ID        uint `gorm:"primaryKey"`
	ClientID  uint
	CarID     uint
	StartTime time.Time
	EndTime   time.Time
	Status    string
	CreatedAt time.Time
}
