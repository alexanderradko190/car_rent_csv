package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"car-export-go/internal/config"
	"car-export-go/internal/service"
	"car-export-go/internal/storage"
)

func main() {
	cfg := config.LoadConfig()
	if err := os.MkdirAll(cfg.ExportDir, 0755); err != nil {
		log.Fatalf("Cannot create export dir: %v", err)
	}
	db, err := storage.NewDB(cfg)
	if err != nil {
		log.Fatalf("Cannot connect DB: %v", err)
	}

	exportService := service.NewExportService(db, cfg.ExportDir)

	http.HandleFunc("/exports/", func(w http.ResponseWriter, r *http.Request) {
        fileName := filepath.Base(r.URL.Path) // Получаем только имя файла (без слэшей)
        filePath := filepath.Join(cfg.ExportDir, fileName)
        w.Header().Set("Content-Type", "text/csv; charset=utf-8")
        w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
        http.ServeFile(w, r, filePath)
    })

	http.HandleFunc("/export/cars", func(w http.ResponseWriter, r *http.Request) {
        path, err := exportService.ExportCars()
        if err != nil {
            http.Error(w, "Ошибка экспорта: "+err.Error(), http.StatusInternalServerError)
            return
        }
        fileName := filepath.Base(path)
        url := "/exports/" + fileName
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":   "done",
            "file_url": url,
        })
    })

	http.HandleFunc("/export/clients", func(w http.ResponseWriter, r *http.Request) {
		path, err := exportService.ExportClients()
		if err != nil {
			http.Error(w, "Ошибка экспорта: "+err.Error(), http.StatusInternalServerError)
			return
		}
		sendExportResult(w, r, path)
	})

	http.HandleFunc("/export/rent_histories", func(w http.ResponseWriter, r *http.Request) {
		path, err := exportService.ExportRentHistories()
		if err != nil {
			http.Error(w, "Ошибка экспорта: "+err.Error(), http.StatusInternalServerError)
			return
		}
		sendExportResult(w, r, path)
	})

	http.HandleFunc("/export/rental_requests", func(w http.ResponseWriter, r *http.Request) {
		path, err := exportService.ExportRentalRequests()
		if err != nil {
			http.Error(w, "Ошибка экспорта: "+err.Error(), http.StatusInternalServerError)
			return
		}
		sendExportResult(w, r, path)
	})

	log.Println("HTTP сервер запущен на :8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}

func sendExportResult(w http.ResponseWriter, r *http.Request, filePath string) {
	fileName := filepath.Base(filePath)
	url := "/exports/" + fileName
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "done",
		"file_url": url,
	})
}
