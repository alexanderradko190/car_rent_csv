package service

import (
	"car-export-go/internal/entity"
	"car-export-go/internal/repository"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type ExportService struct {
	repo      repository.ExportRepository
	exportDir string
}

func NewExportService(repo repository.ExportRepository, exportDir string) *ExportService {
	return &ExportService{repo: repo, exportDir: exportDir}
}

func (s *ExportService) ExportCars() (string, error) {
	cars, err := s.repo.GetCars()
	if err != nil {
		log.Printf("ошибка получения машин: %v", err)
		return "", fmt.Errorf("repo error: %w", err)
	}
	file := s.makeFileName("cars")
	if err := writeCarsCSV(file, cars); err != nil {
		log.Printf("ошибка записи CSV машин: %v", err)
		return "", err
	}
	log.Printf("экспорт машин завершён: %s (кол-во: %d)", file, len(cars))
	return file, nil
}

func (s *ExportService) ExportClients() (string, error) {
	clients, err := s.repo.GetClients()
	if err != nil {
		log.Printf("ошибка получения клиентов: %v", err)
		return "", fmt.Errorf("repo error: %w", err)
	}
	file := s.makeFileName("clients")
	if err := writeClientsCSV(file, clients); err != nil {
		log.Printf("ошибка записи CSV клиентов: %v", err)
		return "", err
	}
	log.Printf("экспорт клиентов завершён: %s (кол-во: %d)", file, len(clients))
	return file, nil
}

func (s *ExportService) ExportRentHistories() (string, error) {
	histories, err := s.repo.GetRentHistories()
	if err != nil {
		log.Printf("ошибка получения историй: %v", err)
		return "", fmt.Errorf("repo error: %w", err)
	}
	file := s.makeFileName("rent_history")
	if err := writeRentHistoriesCSV(file, histories); err != nil {
		log.Printf("ошибка записи CSV истории: %v", err)
		return "", err
	}
	log.Printf("экспорт истории завершён: %s (кол-во: %d)", file, len(histories))
	return file, nil
}

func (s *ExportService) ExportRentalRequests() (string, error) {
	requests, err := s.repo.GetRentalRequests()
	if err != nil {
		log.Printf("ошибка получения заявок: %v", err)
		return "", fmt.Errorf("repo error: %w", err)
	}
	file := s.makeFileName("rental_requests")
	if err := writeRentalRequestsCSV(file, requests); err != nil {
		log.Printf("ошибка записи CSV заявок: %v", err)
		return "", err
	}
	log.Printf("экспорт заявок завершён: %s (кол-во: %d)", file, len(requests))
	return file, nil
}

func (s *ExportService) makeFileName(prefix string) string {
	return filepath.Join(s.exportDir, fmt.Sprintf("%s_%d.csv", prefix, time.Now().Unix()))
}

func writeCarsCSV(filename string, cars []entity.Car) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
    	if err := f.Close(); err != nil {
    		log.Printf("Ошибка при закрытии файла: %v", err)
    	}
    }()
	if _, err := f.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return fmt.Errorf("write BOM: %w", err)
	}
	w := csv.NewWriter(f)
	w.Comma = ';'
	headers := []string{"ID", "Марка", "Модель", "Год", "VIN", "Гос номер", "Класс", "Мощность", "Стоимость за час", "Статус", "ID текущего арендатора"}
	if err := w.Write(headers); err != nil {
		return fmt.Errorf("write headers: %w", err)
	}
	for _, car := range cars {
		currentRenter := ""
		if car.CurrentRenterID.Valid {
			currentRenter = fmt.Sprint(car.CurrentRenterID.Int64)
		}
		row := []string{
			fmt.Sprint(car.ID), car.Make, car.Model, fmt.Sprint(car.Year), car.Vin,
			car.LicensePlate, car.CarClass, fmt.Sprint(car.Power), fmt.Sprintf("%.2f", car.HourlyRate),
			car.Status, currentRenter,
		}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}

func writeClientsCSV(filename string, clients []entity.Client) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
    	if err := f.Close(); err != nil {
    		log.Printf("Ошибка при закрытии файла: %v", err)
    	}
    }()
	if _, err := f.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return fmt.Errorf("write BOM: %w", err)
	}
	w := csv.NewWriter(f)
	w.Comma = ';'
	headers := []string{"ID", "Имя", "Email"}
	if err := w.Write(headers); err != nil {
		return fmt.Errorf("write headers: %w", err)
	}
	for _, c := range clients {
		row := []string{fmt.Sprint(c.ID), c.Name, c.Email}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}

func writeRentHistoriesCSV(filename string, histories []entity.RentHistory) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
    	if err := f.Close(); err != nil {
    		log.Printf("Ошибка при закрытии файла: %v", err)
    	}
    }()
	if _, err := f.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return fmt.Errorf("write BOM: %w", err)
	}
	w := csv.NewWriter(f)
	w.Comma = ';'
	headers := []string{"ID", "ID Клиента", "ID Авто", "Начало аренды", "Завершение аренды", "Общая сумма"}
	if err := w.Write(headers); err != nil {
		return fmt.Errorf("write headers: %w", err)
	}
	for _, h := range histories {
		row := []string{
			fmt.Sprint(h.ID), fmt.Sprint(h.ClientID), fmt.Sprint(h.CarID),
			h.StartTime.Format("2006-01-02 15:04:05"), h.EndTime.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%.2f", h.TotalCost),
		}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}

func writeRentalRequestsCSV(filename string, requests []entity.RentalRequest) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
    	if err := f.Close(); err != nil {
    		log.Printf("Ошибка при закрытии файла: %v", err)
    	}
    }()
	if _, err := f.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		return fmt.Errorf("write BOM: %w", err)
	}
	w := csv.NewWriter(f)
	w.Comma = ';'
	headers := []string{"ID", "ID Клиента", "ID Авто", "Начало аренды", "Завершение аренды", "Статус", "Время создания заявки"}
	if err := w.Write(headers); err != nil {
		return fmt.Errorf("write headers: %w", err)
	}
	for _, r := range requests {
		row := []string{
			fmt.Sprint(r.ID), fmt.Sprint(r.ClientID), fmt.Sprint(r.CarID),
			r.StartTime.Format("2006-01-02 15:04:05"), r.EndTime.Format("2006-01-02 15:04:05"),
			r.Status, r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("write row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}
