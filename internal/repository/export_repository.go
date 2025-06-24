package repository

import "car-export-go/internal/entity"

type ExportRepository interface {
	GetCars() ([]entity.Car, error)
	GetClients() ([]entity.Client, error)
	GetRentHistories() ([]entity.RentHistory, error)
	GetRentalRequests() ([]entity.RentalRequest, error)
}

// Mock реализация (заменить на БД при необходимости)
type mockExportRepository struct{}

func NewExportRepository() ExportRepository {
	return &mockExportRepository{}
}

func (r *mockExportRepository) GetCars() ([]entity.Car, error) {
	return []entity.Car{{1, "Toyota", "Corolla", 2020}}, nil
}

func (r *mockExportRepository) GetClients() ([]entity.Client, error) {
	return []entity.Client{{1, "Ivan", "ivan@mail.com"}}, nil
}

func (r *mockExportRepository) GetRentHistories() ([]entity.RentHistory, error) {
	return []entity.RentHistory{{1, 1, 1, "2024-01-01", "2024-01-03"}}, nil
}

func (r *mockExportRepository) GetRentalRequests() ([]entity.RentalRequest, error) {
	return []entity.RentalRequest{{1, 1, 1, "approved"}}, nil
}
