// Package repository реализует слой доступа к данным (репозитории) для экспорта.
package repository

import (
	"car-export-go/internal/entity"
	"gorm.io/gorm"
)

// ExportRepository определяет интерфейс для экспорта данных из базы.
type ExportRepository interface {
	GetCars() ([]entity.Car, error)
	GetClients() ([]entity.Client, error)
	GetRentHistories() ([]entity.RentHistory, error)
	GetRentalRequests() ([]entity.RentalRequest, error)
}

type exportRepository struct {
	db *gorm.DB
}

// NewExportRepository создаёт новый репозиторий экспорта.
func NewExportRepository(db *gorm.DB) ExportRepository {
	return &exportRepository{db: db}
}

func (r *exportRepository) GetCars() ([]entity.Car, error) {
	var cars []entity.Car
	err := r.db.Find(&cars).Error
	return cars, err
}

func (r *exportRepository) GetClients() ([]entity.Client, error) {
	var clients []entity.Client
	err := r.db.Find(&clients).Error
	return clients, err
}

func (r *exportRepository) GetRentHistories() ([]entity.RentHistory, error) {
	var histories []entity.RentHistory
	err := r.db.Find(&histories).Error
	return histories, err
}

func (r *exportRepository) GetRentalRequests() ([]entity.RentalRequest, error) {
	var requests []entity.RentalRequest
	err := r.db.Find(&requests).Error
	return requests, err
}
