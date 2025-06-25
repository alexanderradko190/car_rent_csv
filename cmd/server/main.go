// Package main запускает HTTP-сервер и слушает очереди RabbitMQ для задач экспорта.
package main

import (
	"car-export-go/internal/config"
	"car-export-go/internal/repository"
	"car-export-go/internal/service"
	"car-export-go/internal/storage"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	cfg := config.LoadConfig()
	db, err := storage.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewExportRepository(db)
	svc := service.NewExportService(repo, cfg.ExportDir)

	go startHTTPServer(svc, cfg.ExportDir)

	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		log.Fatal(err)
	}

	q, err := ch.QueueDeclare("export_tasks", false, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		log.Fatal(err)
	}

	log.Println("Ожидание задач экспорта в RabbitMQ...")

	for d := range msgs {
		task := string(d.Body)
		switch task {
		case "cars":
			if _, err := svc.ExportCars(); err != nil {
				log.Printf("ошибка экспорта cars: %v", err)
			} else {
				log.Println("Экспорт cars завершён")
			}
		case "clients":
			if _, err := svc.ExportClients(); err != nil {
				log.Printf("ошибка экспорта clients: %v", err)
			} else {
				log.Println("Экспорт clients завершён")
			}
		case "rent_histories":
			if _, err := svc.ExportRentHistories(); err != nil {
				log.Printf("ошибка экспорта rent_histories: %v", err)
			} else {
				log.Println("Экспорт rent_histories завершён")
			}
		case "rental_requests":
			if _, err := svc.ExportRentalRequests(); err != nil {
				log.Printf("ошибка экспорта rental_requests: %v", err)
			} else {
				log.Println("Экспорт rental_requests завершён")
			}
		default:
			log.Printf("Неизвестная задача: %s", task)
		}
	}

	_ = ch.Close()
	_ = conn.Close()
}

func startHTTPServer(svc *service.ExportService, exportDir string) {
	http.HandleFunc("/export/cars", func(w http.ResponseWriter, _ *http.Request) {
		path, err := svc.ExportCars()
		if err != nil {
			http.Error(w, "Ошибка экспорта: "+err.Error(), http.StatusInternalServerError)
			return
		}
		sendExportResult(w, path)
	})

	http.HandleFunc("/export/clients", func(w http.ResponseWriter, _ *http.Request) {
		path, err := svc.ExportClients()
		if err != nil {
			http.Error(w, "Ошибка экспорта: "+err.Error(), http.StatusInternalServerError)
			return
		}
		sendExportResult(w, path)
	})

	http.HandleFunc("/export/rent_histories", func(w http.ResponseWriter, _ *http.Request) {
		path, err := svc.ExportRentHistories()
		if err != nil {
			http.Error(w, "Ошибка экспорта: "+err.Error(), http.StatusInternalServerError)
			return
		}
		sendExportResult(w, path)
	})

	http.HandleFunc("/export/rental_requests", func(w http.ResponseWriter, _ *http.Request) {
		path, err := svc.ExportRentalRequests()
		if err != nil {
			http.Error(w, "Ошибка экспорта: "+err.Error(), http.StatusInternalServerError)
			return
		}
		sendExportResult(w, path)
	})

	http.HandleFunc("/exports/", func(w http.ResponseWriter, r *http.Request) {
		fileName := filepath.Base(r.URL.Path)
		filePath := filepath.Join(exportDir, fileName)
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		http.ServeFile(w, r, filePath)
	})

	log.Println("HTTP сервер запущен на :8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}

func sendExportResult(w http.ResponseWriter, filePath string) {
	fileName := filepath.Base(filePath)
	url := "/exports/" + fileName
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "done",
		"file_url": url,
	}); err != nil {
		log.Printf("Ошибка при кодировании ответа: %v", err)
	}
}
