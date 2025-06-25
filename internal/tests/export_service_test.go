package tests

import (
	"car-export-go/internal/entity"
	"car-export-go/internal/service"
	"errors"
	"os"
	"testing"
)

type fakeRepo struct{}

func (f *fakeRepo) GetCars() ([]entity.Car, error) { return []entity.Car{{ID: 1, Make: "BMW"}}, nil }
func (f *fakeRepo) GetClients() ([]entity.Client, error) { return []entity.Client{{ID: 1, Name: "Ivan"}}, nil }
func (f *fakeRepo) GetRentHistories() ([]entity.RentHistory, error) { return nil, errors.New("fail") }
func (f *fakeRepo) GetRentalRequests() ([]entity.RentalRequest, error) { return []entity.RentalRequest{}, nil }

func TestExportCars(t *testing.T) {
    t.Log("-> стартует TestExportCars")
	dir := os.TempDir()
	s := service.NewExportService(&fakeRepo{}, dir)
	file, err := s.ExportCars()
	if err != nil {
		t.Fatalf("ошибка экспорта: %v", err)
	}
	if _, err := os.Stat(file); err != nil {
		t.Fatalf("файл не найден: %v", err)
	}
	os.Remove(file)
}

func TestExportRentHistories_Error(t *testing.T) {
    t.Log("-> стартует TestExportRentHistories_Error")
	dir := os.TempDir()
	s := service.NewExportService(&fakeRepo{}, dir)
	_, err := s.ExportRentHistories()
	if err == nil {
		t.Fatal("ожидалась ошибка экспорта")
	}
}
