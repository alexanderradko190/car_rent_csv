package service

import (
    "car-export-go/internal/entity"
    "encoding/csv"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "time"
    "gorm.io/gorm"
)

type ExportService struct {
    db        *gorm.DB
    exportDir string
}

func NewExportService(db *gorm.DB, exportDir string) *ExportService {
    return &ExportService{
        db:        db,
        exportDir: exportDir,
    }
}

func (s *ExportService) ExportCars() (string, error) {
    var cars []entity.Car
    if err := s.db.Find(&cars).Error; err != nil {
        return "", err
    }
    filename := filepath.Join(s.exportDir, "cars_"+strconv.FormatInt(time.Now().Unix(), 10)+".csv")
    f, err := os.Create(filename)
    if err != nil {
        return "", err
    }
    defer f.Close()

    bom := []byte{0xEF, 0xBB, 0xBF}
    if _, err := f.Write(bom); err != nil {
        return "", err
    }

    writer := csv.NewWriter(f)
    writer.Comma = ';'

    headers := []string{
        "ID",
        "Марка",
        "Модель",
        "Год",
        "VIN",
        "Гос номер",
        "Класс",
        "Мощность",
        "Стоимость за час",
        "Статус",
        "ID текущего арендатора",
    }
    if err := writer.Write(headers); err != nil {
        return "", err
    }

    for _, car := range cars {
        currentRenter := ""
        if car.CurrentRenterID.Valid {
            currentRenter = fmt.Sprint(car.CurrentRenterID.Int64)
        }
        row := []string{
            fmt.Sprint(car.ID),
            car.Make,
            car.Model,
            fmt.Sprint(car.Year),
            car.Vin,
            car.LicensePlate,
            car.CarClass,
            fmt.Sprint(car.Power),
            fmt.Sprintf("%.2f", car.HourlyRate),
            car.Status,
            currentRenter,
        }
        if err := writer.Write(row); err != nil {
            return "", err
        }
    }
    writer.Flush()
    return filename, writer.Error()
}

func (s *ExportService) ExportClients() (string, error) {
    var clients []entity.Client
    if err := s.db.Find(&clients).Error; err != nil {
        return "", err
    }
    filename := filepath.Join(s.exportDir, "clients_"+strconv.FormatInt(time.Now().Unix(), 10)+".csv")
    f, err := os.Create(filename)
    if err != nil {
        return "", err
    }
    defer f.Close()

    bom := []byte{0xEF, 0xBB, 0xBF}
    if _, err := f.Write(bom); err != nil {
        return "", err
    }

    writer := csv.NewWriter(f)
    writer.Comma = ';'

    headers := []string{
        "ID",
        "Имя",
        "Email",
    }
    if err := writer.Write(headers); err != nil {
        return "", err
    }
    for _, c := range clients {
        row := []string{
            fmt.Sprint(c.ID),
            c.Name,
            c.Email,
        }
        if err := writer.Write(row); err != nil {
            return "", err
        }
    }
    writer.Flush()
    return filename, writer.Error()
}

func (s *ExportService) ExportRentHistories() (string, error) {
    var histories []entity.RentHistory
    if err := s.db.Find(&histories).Error; err != nil {
        return "", err
    }
    filename := filepath.Join(s.exportDir, "rent_history_"+strconv.FormatInt(time.Now().Unix(), 10)+".csv")
    f, err := os.Create(filename)
    if err != nil {
        return "", err
    }
    defer f.Close()

    bom := []byte{0xEF, 0xBB, 0xBF}
    if _, err := f.Write(bom); err != nil {
        return "", err
    }

    writer := csv.NewWriter(f)
    writer.Comma = ';'

    headers := []string{
        "ID",
        "ID Клиента",
        "ID Авто",
        "Начало аренды",
        "Завершение аренды",
        "Общая сумма",
    }
    if err := writer.Write(headers); err != nil {
        return "", err
    }
    for _, h := range histories {
        row := []string{
            fmt.Sprint(h.ID),
            fmt.Sprint(h.ClientID),
            fmt.Sprint(h.CarID),
            h.StartTime.Format("2006-01-02 15:04:05"),
            h.EndTime.Format("2006-01-02 15:04:05"),
            fmt.Sprintf("%.2f", h.TotalCost),
        }
        if err := writer.Write(row); err != nil {
            return "", err
        }
    }
    writer.Flush()
    return filename, writer.Error()
}

func (s *ExportService) ExportRentalRequests() (string, error) {
    var requests []entity.RentalRequest
    if err := s.db.Find(&requests).Error; err != nil {
        return "", err
    }
    filename := filepath.Join(s.exportDir, "rental_requests_"+strconv.FormatInt(time.Now().Unix(), 10)+".csv")
    f, err := os.Create(filename)
    if err != nil {
        return "", err
    }
    defer f.Close()

    bom := []byte{0xEF, 0xBB, 0xBF}
    if _, err := f.Write(bom); err != nil {
        return "", err
    }

    writer := csv.NewWriter(f)
    writer.Comma = ';'

    headers := []string{
        "ID",
        "ID Клиента",
        "ID Авто",
        "Начало аренды",
        "Завершение аренды",
        "Статус",
        "Время создания заявки",
    }
    if err := writer.Write(headers); err != nil {
        return "", err
    }
    for _, r := range requests {
        row := []string{
            fmt.Sprint(r.ID),
            fmt.Sprint(r.ClientID),
            fmt.Sprint(r.CarID),
            r.StartTime.Format("2006-01-02 15:04:05"),
            r.EndTime.Format("2006-01-02 15:04:05"),
            r.Status,
            r.CreatedAt.Format("2006-01-02 15:04:05"),
        }
        if err := writer.Write(row); err != nil {
            return "", err
        }
    }
    writer.Flush()
    return filename, writer.Error()
}

