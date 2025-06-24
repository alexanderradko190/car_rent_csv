package entity

import (
	"database/sql" // <--- Это обязательно!
	"time"
)

type Car struct {
	ID              uint         `gorm:"primaryKey"`
	Make            string
	Model           string
	Year            int
	Vin             string
	LicensePlate    string
	CarClass        string
	Power           int
	HourlyRate      float64
	Status          string
	CurrentRenterID sql.NullInt64 // <-- теперь Go точно знает, что это за тип
}

type Client struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string
	Email string
}

type RentHistory struct {
	ID        uint      `gorm:"primaryKey"`
	ClientID  uint
	CarID     uint
	StartTime time.Time
	EndTime   time.Time
	TotalCost float64
}

type RentalRequest struct {
	ID        uint      `gorm:"primaryKey"`
	ClientID  uint
	CarID     uint
	StartTime time.Time
	EndTime   time.Time
	Status    string
	CreatedAt time.Time
}
